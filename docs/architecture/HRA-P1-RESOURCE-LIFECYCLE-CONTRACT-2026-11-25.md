# HRA-P1 — Resource Lifecycle Phase-One Contract

## Contract paper `HRA-P1-2026-11-25-01`

| Field | Value |
|---|---|
| Date | **2026-11-25** |
| Type | Contract paper — verified research draft |
| Authority | Source audit of the resource, image, chat, audio, workflow, and process paths |
| Gate | Phase 1 — lifecycle, accounting, enforcement, recovery |
| Status | **PROPOSED — characterization complete; Gate 1 implemented** |
| Implementation authorization | **Gate 1 implemented; later gates pending** |
| Baseline tests | Targeted suites pass; see §14 |
| User-visible workflow demonstrated | **No** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Proposed** · **Defect** · **Unestablished**.

**Identifier note:** This paper is self-contained. Clause identifiers are stable within this paper; later amendments must add clauses rather than silently rewrite the meaning of an existing clause.

---

## 1. Contract decision

**RLC-01 (Proposed — scope).** This contract defines phase one of Howl's resource lifecycle for chat, image, video, ASR, TTS, and sequential image workflows. Phase one consists of four properties: **ownership**, **accounting**, **enforcement**, and **recovery**.

**RLC-02 (Proposed — one authority, three views).** One resource authority owns admission, ownership leases, capacity commitments, and the transitions between them. It exposes three distinct views:

| View | Owns | Does not own |
|---|---|---|
| Ownership | Leases, owners, process identity, resident registration, workflow claims | Physical memory readings |
| Capacity | Outstanding commitments, pending deductions, reflected allocations, recovery | Process lifetime by itself |
| Observation | Physical free/total, managed placement, measured peaks, confidence, age | Permission to allocate |

These are views of one authority, not three independent ledgers.

**RLC-03 (Proposed — out of scope).** Persistent calibration, learned placement optimization, dynamic reserve optimization, UI design, and changes to model behavior are outside this contract. They begin only after phase one passes.

**RLC-04 (Proposed — acceptance principle).** A successful operation without a correctly closed lease is not a successful operation.

**RLC-05 (Proposed — existing behavior is preserved unless explicitly changed).** Image, chat, audio, and video behavior remains intact unless this contract, or a later authorized change, explicitly replaces a behavior and proves the replacement.

---

## 2. Current owners and competing owners

### 2.1 Current owners

| State | Current owner | Evidence |
|---|---|---|
| Admission decision | `resource.Authority` | `internal/resource/admission.go` |
| Loaded-model registry | `resource.ResourceManager` | `internal/resource/manager.go` |
| Image registration | `App.registerImageModel` / `unregisterImageModels` | `app-image-model.go:1015`, `:1026` |
| Image process | `sdcpp.Engine.serverProc` | `internal/image/engine/sdcpp/engine.go` |
| Video process | `sdcpp.Engine.videoProc` | `internal/image/engine/sdcpp/video.go` |
| Chat process | `host.Processes` | `internal/host/processes.go` |
| Chat registration | load monitor | `app-model.go:524` |
| Audio registration | `Pipeline.registerAudio` | `internal/audio/pipeline.go:704` |
| Image reclamation | `App.imageReclaim` | `app-image-reclaim.go` |
| Workflow sequencing | `App.imageStartMu` and `ImageWorkflowStage` | `app-image-workflow-run.go` |
| Separate allocation advice | `resource.Allocate` / `RequestAllocation` | `internal/resource/allocator.go` |

**Missing owner (Observed).** Chat execution has no resource-layer execution lease; the streaming client is not an authority-owned scope. Audio has an audio-manager execution lease, but its resource lease is released at registration. These are gaps in the lifecycle contract, not alternative ownership models.

### 2.2 Competing owners found

| Competing owner | Conflict |
|---|---|
| `Engine.Shutdown` unregisters all image models | Duplicates `App.unregisterImageModels` and bypasses any lease transition |
| `App.UnloadImageModel` unregisters after `Shutdown` | Duplicate cleanup with different ownership assumptions |
| Public `App.RegisterLoadedModel` / `UnregisterLoadedModel` | Direct registry mutation with no lease or process identity |
| `videoProc` and `serverProc` | Two process owners outside `host.Processes` |
| Job manager `onFinished` callback | Reclamation lifetime is not tied to an execution lease |
| `Image.Reconfigure` / `SetEngine` | Can replace the engine under a workflow without a workflow-visible transition |
| `resource.Allocate` | Separate 90%/95% policy that can disagree with the authority |
| `imageReclaim` reserve | Independent 10% policy, separate from authority reserve |

### 2.3 What stays

- `AdmissionRequest`, `AdmissionEstimate`, `AdmissionDecision` remain the decision vocabulary.
- `AdmissionLease` remains the capacity reservation handle, but its meaning is narrowed to **reservation only**.
- `ResourceManager` remains the resident registry view, but it becomes authority-written rather than caller-written.
- `sdcpp.Engine`, llama server, audio loader, and job manager remain execution adapters.
- `RunImageWorkflow` remains the sequential workflow owner, but receives a workflow scope and component claims.

### 2.4 What is removed or replaced

| Removed / competing path | Replacement |
|---|---|
| `Engine.Shutdown` directly unregistering models | Engine reports process exit; authority performs the registration transition |
| `App.UnloadImageModel` directly unregistering | Authority-owned unload transition |
| Public registration bindings as an authority bypass | Removed from the public surface or restricted to an authority-owned adapter |
| `host.Processes` single-child ownership | One process registry with multiple roles, preserving llama behavior |
| `video.stopServer` direct `Wait()` | Registry/reaper-owned process exit |
| `resource.Allocate` as a second policy | Projection from the authority |
| `imageReclaim` as a separate reserve authority | Authority-provided facts and policy |

---

## 3. Authorities and boundaries

**RLC-06 (Proposed — authority table).**

| Authority | Owns | Must not own |
|---|---|---|
| Resource Authority | Admission, leases, capacity commitments, transition ordering | Model semantics, chat history, UI copy |
| Process Registry | Child identity, start token, single `Wait`, exit observation, stop | Whether a workload is admitted |
| Engine Adapter | Start, readiness, execution, measured evidence, exit report | Global admission or registration mutation |
| Observation Layer | Physical and managed measurements, peaks, age, confidence | Permission beyond current telemetry and limits |
| Workflow Scope | Parent lifetime, stage order, retained claims, cancellation | A second admission policy |
| Configuration Authority | User intent, manual limits, provenance | Physical capacity or process lifetime |

**RLC-07 (Proposed — no inferred permission).** No component may treat its own estimate, cached model state, or historical measurement as permission from another component. Physical observation and the authority's current commitments decide admission.

---

## 4. Invariants

These are the acceptance boundary. A type or test that does not enforce one of these is not phase-one work.

### 4.1 Ownership

**RLC-08 — One registration owner.** Every model registration has exactly one authoritative creator and one authoritative remover. Direct callers may request a transition; they may not write the registry.

**RLC-09 — Registration is process-backed.** A registration may exist only while its associated process identity is alive and its readiness evidence is current. No registration may outlive its process.

**RLC-10 — One child, one `Wait`.** Every child process has exactly one owner of `Wait()`. A second waiter is a defect, even if it usually returns safely.

**RLC-11 — Every process-backed workload has an execution scope.** Chat generation, image generation, video rendering, audio processing, and every workflow stage have an identifiable execution scope owned by the authority. A scope that consumes resource capacity carries an execution lease. A scope that cannot consume capacity must record why it has no lease; absence of a lease is never left implicit.

**RLC-11a — Execution scopes are named, not inferred.** At minimum: one chat turn or stream, one image job, one video render, one audio request, and one workflow stage. A generic "engine is busy" flag is not an execution scope, because it does not identify the owner, the resource commitment, or the cleanup transition.

**RLC-12 — Component ownership is explicit.** A retained component has a stable component key, an owner, and a reference count. Reuse is a claim, not an assumption.

**RLC-13 — Engine identity is part of ownership.** A workflow or execution scope names the engine instance it uses. Replacing or stopping that engine invalidates the scope; it cannot continue silently on a stale pointer.

### 4.2 Capacity

**RLC-14 — Incremental accounting.** Physical free memory already includes resident allocations. The authority charges only outstanding/unreflected commitments and the new workload's incremental peak; it must not charge resident weights again.

**RLC-15 — Pending is an unseen commitment.** A capacity commitment is `pending` until physical telemetry reflects the allocation. While pending, it is deducted from effective free capacity. Once reflected, it stops being an additional deduction because physical free now includes it.

**RLC-15a — What "reflected" means.** A commitment is reflected only when a post-start observation accounts for the allocation: free memory has moved in the expected direction by at least the committed amount, or the allocation is explicitly attributed by managed or process evidence. Exact equality is not required. If the observation is ambiguous, incomplete, or cannot distinguish this allocation from other activity, the commitment remains `pending` or becomes `unknown`; it is not silently released.

**RLC-16 — Reflected is not released.** A commitment becoming reflected does not release its ownership lease. The lease remains active until its lifecycle owner ends and the registration/process transition completes.

**RLC-17 — Capacity recovery is a separate event.** Process exit invalidates ownership, but recovered capacity is established only by a subsequent physical measurement (or an explicit unknown state). A killed process is not proof that the driver returned the memory.

**RLC-18 — Physical free is the hard gate.** Managed placement and engine budgets are not proof of physical headroom. `--max-vram` is a managed budget, not a physical cap.

**RLC-19 — No double counting.** A replacement credit is conditional and is committed only after observed release. The authority must not add a resident estimate to free memory that already excludes it.

**RLC-20 — RAM is a first-class capacity.** RAM-only residents participate in blocking and accounting; `GPUMemoryMB <= 0` is not invisibility.

### 4.3 Enforcement

**RLC-21 — Atomic transitions.** A successful load cannot become resident without a resident lease. A failed transfer leaves no dual lease. Execution cannot outlive its stage or process. Reclamation requires no active owner.

**RLC-22 — Shared components are reference-counted.** A stage acquiring an already-resident component increments a claim; releasing it decrements. Reclamation is refused while the count is non-zero.

**RLC-23 — Workflow scope reserves only retained and explicitly committed capacity.** It does not reserve every future stage. A next stage's capacity is committed when the workflow explicitly advances to it.

**RLC-24 — Manual limits preserve provenance.** Manual, automatic, clamped-manual, zero/live-free, and negative-reserve values remain distinct across reloads and workflows.

**RLC-25 — Optional telemetry is not a gate.** Missing CPU utilization, temperature, or per-process VRAM does not by itself refuse work. Only physical memory telemetry, explicit safety rules, and current manual limits may gate admission.

**RLC-26 — Historical observations never grant permission.** They improve prediction only. Underprediction may tighten immediately; reducing safety margins requires stronger current evidence.

### 4.4 Recovery

**RLC-27 — Cleanup is idempotent.** Cancellation, shutdown, callback completion, process exit, and reclamation may race; each transition closes at most once and never leaves a phantom reservation or registration.

**RLC-28 — No orphan process or lease.** Shutdown closes admission, joins execution scopes, stops owned processes, observes exits, and reconciles capacity. A process or lease that outlives shutdown is a recovery failure.

**RLC-29 — Refusals are surfaced.** A discarded refusal is non-conforming.

---

## 5. State model

### 5.1 Lease kinds

| Kind | Purpose | Lifetime |
|---|---|---|
| Load | Reserve capacity for a start that has not yet become resident | From admission until successful commit or failure |
| Resident | Own a loaded model and its registration | From committed readiness until process exit and registration removal |
| Execution | Own one generation, stage, render, or audio request | From stage start until execution settles |
| Workflow | Own a multi-stage parent and its retained claims | From workflow open until close, failure, or cancellation |

A lease has an owner, a scope, and a creation identity. `AdmissionLease` is the capacity-reservation projection of the load lease; it is not a generic release for any lifecycle event.

### 5.2 Capacity states

| State | Meaning | Deducted from effective free? | Ownership still active? |
|---|---|---|---|
| `pending` | Committed but not yet visible in physical telemetry | Yes | Yes |
| `reflected` | A post-start observation accounts for the allocation | No | Yes |
| `released` | Allocation has ended and recovery is established | No | No |
| `unknown` | Exit or observation could not establish recovery | Yes, conservatively | No or invalidated |

`pending → reflected` is a capacity transition only. It does not close a resident or execution lease.

### 5.3 Process and registration states

| Object | States |
|---|---|
| Process | `starting` → `alive` → `exiting` → `exited` → `reaped`; any failure may produce `unknown` |
| Registration | `none` → `pending` → `resident`; `resident → invalid` on process exit; `invalid → removed` after the authority transition |
| Engine instance | `starting` → `ready`; `ready → replaced` or `ready → stopped` |
| Workflow | `open` → `stage_active` → `between_stages` → `closing` → `closed`; failure may produce `failed` |

A registration in `pending` may not be treated as resident. A registration in `invalid` may not be used for admission.

### 5.4 Required transitions

| Transition | Owner | Required result |
|---|---|---|
| Request → admit | Resource Authority | Capacity lease granted or explicit refusal; no process started |
| Admit → start | Engine Adapter + Process Registry | Process identity allocated; load lease active; registration pending |
| Ready | Engine Adapter | Process alive; readiness evidence; registration resident; load lease committed |
| Observe allocation | Observation Layer | Capacity `pending → reflected`; ownership unchanged |
| Start execution | Resource Authority + Engine Adapter | Execution lease active; incremental capacity pending |
| Observe execution | Observation Layer | Execution capacity `pending → reflected`; execution lease active |
| Settle execution | Execution owner | Execution lease released exactly once; resident lease remains if policy retains the model |
| Retain component | Workflow Scope | Component claim incremented with stable key |
| Release component | Workflow Scope | Claim decremented; reclamation refused while non-zero |
| Process exit | Process Registry | Process `exited`; registration `invalid`; resident/execution leases invalidated |
| Recover capacity | Observation Layer + Authority | Capacity `released` only after physical measurement; otherwise `unknown` |
| Unload | Resource Authority | Drain, stop, observe exit, invalidate registration, recover capacity |
| Workflow close | Workflow Scope | Stage leases settled; claims released; parent closed |
| Shutdown | Lifecycle Authority | Admission closed; execution scopes joined; processes stopped; capacity reconciled |

### 5.5 The three distinct release events

**RLC-30 (Proposed — release vocabulary).**

| Event | What ends | What does not end |
|---|---|---|
| Reservation release | A load or execution capacity commitment whose operation ended without a resident ownership transition | Resident ownership; process lifetime; physical recovery evidence |
| Ownership release | A resident or workflow lease after its process/registration transition | Physical capacity recovery evidence |
| Capacity recovery | The capacity deduction after a confirmed exit measurement | Historical evidence or the fact that the process existed |

These events must have separate APIs or explicitly distinct transition names. A single `Release()` must not imply all three.

---

## 6. Ownership rules

### 6.1 Registration ownership

**RLC-31.** The authority is the only writer of resident registrations. Engine adapters report:

- process identity;
- readiness evidence;
- actual placement/measurement;
- process exit.

They do not unregister directly.

**RLC-32.** The registration key includes the workload kind, authority owner identity, process key, and model/component identity. Name-only registration is insufficient for shared engines and repeated model names.

**RLC-33.** `RegisterLoadedModel` and `UnregisterLoadedModel` are not authority APIs. They are removed from the public surface or converted into requests routed to the authority.

### 6.2 Process ownership

**RLC-34.** The process registry owns every child process: `llama-server`, `sd-server`, and `sd-cli`. The existing llama behavior is preserved: one recorded handle, one waiter, a stop that reports failure, and ownership-safe clearing.

**RLC-35.** The registry key is role + PID + start token. A reused PID is not a reused process identity.

**RLC-36.** Video `sd-cli` is registered before `Start`, assigned to the owning job before resume, and removed only after `Wait` returns. The existing `videoProc`/`videoOwner` fields may remain as adapter state, but they are not the authority.

**RLC-37.** `stopServer` must not call `Wait()`. It requests a registry stop and awaits the registry's exit signal. A second `Wait()` is a contract violation.

### 6.3 Engine identity

**RLC-38.** Each engine instance has an identity that includes its process identity and startup configuration identity. A workflow or execution lease references that identity.

**RLC-39.** `Reconfigure`, `UnloadImageModel`, and any replacement transition invalidate the affected engine identity before the new engine is visible. An active workflow receives an explicit invalidated result; it does not continue on a stale pointer.

### 6.4 Workflow ownership

**RLC-40.** A workflow has one parent lease. Its stages have execution leases. Retained components have claims under the parent.

**RLC-41.** Component key is stable across stages: engine identity + component role + normalized path + weight type + companion set. A file-size or display name alone is not identity.

**RLC-42.** Stage 2's upscaler is retained only by an explicit claim. Stage 4 must verify the claim; if it is absent, the stage loads the component or refuses. It may not assume residency from the workflow's history.

**RLC-43.** A workflow batch is one parent scope; each item is a separate execution scope. Reclamation at an item boundary may preserve resident components but must not release the parent.

---

## 7. Accounting rules

### 7.1 Effective free capacity

For a target device:

```text
effective_free = max(0, physical_free - outstanding_pending_commitments)
```

- `outstanding_pending_commitments` includes load and execution commitments whose `capacity_state` is `pending` or `unknown`.
- A commitment with `capacity_state = reflected` is no longer deducted; physical free already reflects it.
- A planned replacement credit is not added until its process exit and recovery measurement are observed.
- Resident allocations are not subtracted again; physical free already includes them.
- Reserve is computed from effective free and/or explicit policy, then subtracted before admission.
- Physical free is still the hard gate after managed placement is planned.

### 7.2 Incremental execution peak

**RLC-44.** For an already-resident model, execution reserves only the additional peak required by the execution:

```text
execution_increment =
    max(0, total_workload_peak - resident_footprint_already_in_physical_free)
```

`total_workload_peak` is the total process peak including resident components, new adapters, activation, and fixed overhead. It is not an already-incremental estimate. If an adapter supplies an incremental `ProjectedPeakVRAMMB` for a resident stage, it must be normalized to the total before applying this formula; subtracting a resident footprint from an already-incremental value would under-reserve.

For a workflow stage, `resident_footprint_already_in_physical_free` includes the base model and any claimed retained components. New adapters, activation, and fixed overhead are charged; resident weights are not.

**RLC-45.** For a full load, the load reservation charges the total projected resident + activation + overhead that admission judged, not merely `RequiredVRAMMB` when that is smaller.

**RLC-46.** A `pending` execution commitment remains deducted until a post-start observation reflects it. If observation fails, it remains conservatively `unknown`; it is not silently released.

### 7.3 Predicted, observed, and complete

**RLC-47.** Keep these separate:

| Fact | Meaning | May grant permission? |
|---|---|---|
| Predicted peak | Planning estimate | No |
| Highest observed sample | Lower bound on actual use | No |
| Measurement completeness/confidence | What the sample interval and source prove | No |
| Current physical free | Current capacity evidence | Yes, under policy |
| Manual limit | User/policy constraint | Yes, as a ceiling |

**RLC-48.** Historical observations improve prediction only. They never override current physical telemetry or manual limits.

---

## 8. Manual limits and provenance

**RLC-49.** Budget source is explicit: `manual`, `automatic`, `manual-clamped`, `live-free`, `negative-reserve`, or `unset`.

**RLC-50.** A manual limit change applies to new leases. Existing execution continues unless it crosses a hard safety boundary. An over-limit resident state is reported, not silently reclassified as newly admitted.

**RLC-51.** Disruptive replanning (unload, replace, or restart) requires explicit authorization. A limit change alone is not authorization to terminate active work.

**RLC-52.** `MaxVRAM = 0`, negative reserve, and clamped values preserve their original intent. A clamp to the automatic cap must not erase manual provenance.

---

## 9. Observation and telemetry boundary

**RLC-53.** Physical VRAM/RAM telemetry is the admission signal. CPU utilization, temperature, and per-process VRAM are observations until a specific safety rule says otherwise.

**RLC-54.** Missing optional telemetry lowers confidence or leaves an optional rule unevaluated; it does not refuse by default.

**RLC-55.** Stale, partial, or contradictory physical telemetry fails conservatively. The reason names the source and the missing evidence.

**RLC-56.** A free-memory increase is not capacity evidence until the process exit or authoritative measurement confirms release.

**RLC-57.** Persistent calibration is prohibited until ownership, accounting, enforcement, and recovery pass the phase-one gate.

---

## 10. Failure and recovery matrix

| Failure | Ownership result | Capacity result | Required user-visible result |
|---|---|---|---|
| Admission refused | No lease, no process | No change | Exact refusal with blocker and numbers |
| Load fails before process | Load lease released | Pending removed | Explicit failure; no registration |
| Load fails after old process stopped | Old registration invalidated | Old capacity `unknown` until measured | Explicit failure; no phantom model |
| Process crashes during execution | Execution and resident invalidated | Capacity `unknown` | Execution failure; no ready claim |
| Cancellation | Execution released; resident per policy | Reflected capacity remains | Cancelled result; no phantom execution |
| Unload | Resident invalidated; process stopped | Recovered only after measurement | Unloaded state only after evidence |
| Shutdown | All leases closed; all processes stopped | Reconcile or report unknown | Clean shutdown or named failure |
| Stale telemetry | No new lease | Existing execution continues unless hard boundary | Refusal names stale source |
| Manual limit change | New leases use new limit | Existing commitments remain | Report over-limit resident state |
| Workflow cancellation between stages | Parent closes; claims release | Stage capacity settles | No next stage starts |
| Workflow cancellation during an active stage | The active generation observes the workflow context; parent and claims close once | Engine joins the generation and stops/reaps the owned process | Cancellation reports only after termination, or reports the termination failure |
| Component lost between stages | Stage refuses or reloads | New component charge | No silent reuse claim |

---

## 11. Characterization tests

These establish current behavior before the new types change it. The passing tests document semantics that must survive the migration. The former known-failure tests have graduated into the default suite; no tagged known-failure test remains. The planned acceptance matrix in §12 remains outstanding.

| Test | Package | What it establishes | Status |
|---|---|---|---|
| `TestAdmissionReleaseIsReservationOnly` | `internal/resource` | `Release()` removes pending capacity and does not unregister | Implemented — passes |
| `TestPendingReservationBlocksSecondAdmission` | `internal/resource` | An open reservation prevents a second request from spending the same capacity | Implemented — passes |
| `TestResidentStageIsNotChargedResidentWeightsAgain` | `internal/resource` | A resident stage is judged on its incremental peak, not resident weights again | Implemented — passes |
| `TestGrantReservesProjectedPeak` | `internal/resource` | `grant` records the need admission judged, not only `RequiredVRAMMB` | Graduated — passes (Gate 1) |
| `TestPendingCommitmentSurvivesUntilObservation` | `internal/resource` | Pending deduction remains until a post-start snapshot; reflected does not release ownership | Planned — requires an observation transition |
| `TestRAMOnlyResidentBlocksRAMAdmission` | `internal/resource` | `GPUMemoryMB <= 0` does not make a resident invisible | Graduated — passes (Gate 1) |
| `TestStaleChatProcessCannotRegisterAfterStop` / `TestCurrentChatProcessRegistersAfterHealth` | `app` | A health result registers only for the process the load still owns | Graduated — passes (Gate 2 behavioral) |
| `TestImageReloadFailureInvalidatesPreviousRegistration` / `TestImageReloadSuccessLeavesExactlyOneRegistration` | `app` | A failed replacement leaves no registration; a successful replacement leaves exactly one | Graduated — passes (D-4 behavioral) |
| `TestShutdownJoinsGeneration` | `sdcpp` | Shutdown waits for the active generation scope | Graduated — passes (D-5 behavioral) |
| `TestEngineShutdownReportsExitWithoutUnregistering` | `sdcpp` | The engine reports process exit; the app/domain owns registration invalidation | Graduated — passes (D-5 behavioral) |
| `TestStopServerUsesOneRegisteredWaiter` / `TestStopServerRetainsOwnershipWhenRegisteredStopFails` | `sdcpp` | `stopServer` uses the registry waiter once and retains ownership when exit is not observed | Graduated — passes (Gate 2 behavioral) |
| `TestVideoExecutionHasLease` | `internal/image/job` | Video admission does not end before `sd-cli` runs | Graduated — passes (D-7 behavioral) |
| `TestVideoLeavesNoImageRegistration` | `sdcpp` / `app` | Stopping `sd-server` for video removes/invalidates its registration | Graduated — passes (D-7 behavioral) |
| `TestWorkflowReusesEngineAcrossStages` | `internal/resource` / `app` | Normal same-model workflow keeps one engine instance and rejects a replacement | Graduated — passes (D-8 behavioral) |
| `TestWorkflowRetainedUpscalerIsClaimed` | `internal/resource` / `app` | Stage 2 retains the upscaler; stage 4 verifies the claim | Graduated — passes (D-8 behavioral) |
| `TestOptionalTelemetryAbsenceDoesNotRefuse` | `internal/resource` | Missing CPU/temp/per-process data is not a default refusal | Planned — no optional fields yet |
| `TestManualBudgetProvenanceSurvivesReload` | `app` | Manual, clamped, zero, and negative cases remain distinct | Planned — existing tests are partial |
| `TestPublicRegistrationBypassesAuthority` | `app` | The public binding is either removed or routed through the authority | Planned — design decision |

### 11.1 Known-failure retirement

The default suite must stay green. The Phase-one source guards and tagged known-failure tests have been retired; the corresponding behavioral tests now run with the default packages. The planned acceptance matrix in §12 remains outstanding.

## 12. Phase-one acceptance tests

These must pass before phase one is complete. They are not optional polish.

| Test | Proves |
|---|---|
| `TestConcurrentAdmissionCannotSpendLastVRAM` | Atomic admission and pending commitments |
| `TestPendingExecutionChargedUntilReflected` | Pending/reflected distinction |
| `TestReflectedExecutionKeepsLeaseActive` | Capacity visibility does not release ownership |
| `TestLastProcessExitDoesNotFreeUntilMeasured` | Exit and recovery are separate |
| `TestResidentExecutionDoesNotRechargeWeights` | Incremental accounting |
| `TestSharedComponentRefcountAcrossStages` | Workflow component claims |
| `TestWorkflowCancelBetweenStages` | Parent close and no next stage |
| `TestManualLimitChangeDuringExecution` | New leases use new limit; existing work continues |
| `TestShutdownJoinsAllExecutionScopes` | No callback or process outlives shutdown |
| `TestRegistrationCannotOutliveProcess` | Process/registration binding |
| `TestVideoProcessRegisteredAndReleased` | Video ownership |
| `TestVideoSingleWaitOwner` | Process registry invariant |
| `TestOptionalTelemetryAbsenceDoesNotRefuse` | Optional observation boundary |
| `TestPredictedObservedConfidenceRemainSeparate` | Evidence model |
| `TestCrossWorkloadLeaseSequence` | Chat resident, image executing, video admitted, image completes, video starts, chat executes, all close |

### 12.1 Cross-workload acceptance scenario

The test must use a stateful fake runtime with physical free, free RAM, resident registrations, process identities, and pending/reflected commitments. The expected result is calculated from those values, not from workload labels.

Sequence:

1. Chat model remains resident.
2. Image generation acquires an execution lease and begins.
3. Video requests admission while image execution is active.
4. Image completes and retains its model.
5. Video starts.
6. Chat executes while video remains active, if measured budgets permit.
7. All processes and leases terminate.

Assertions:

- Image execution charges only its incremental peak; the resident chat model is not charged again.
- Video's decision uses physical free minus outstanding pending commitments and its own incremental peak.
- Video is not admitted or refused because of workload labels or arbitrary priority.
- Image completion releases its execution lease but not its resident lease.
- Chat execution is admitted only if its incremental execution lease fits.
- Every lease closes exactly once.
- No process remains alive.
- No registration remains without a live process.
- No pending commitment survives the final reconciled measurement.

### 12.2 Cross-boundary baseline (frozen 2026-11-25)

The D-7 and D-8 named tests are behaviorally verified. Workflow cancellation now has a backend execution-scope implementation: every workflow stage generation is a child of one cancellable context, stage boundaries refuse to start after cancellation, `CancelImageWorkflow` waits for execution cleanup, and accepted cancellation joins the active engine generation and reaps the owned process through `Engine.Shutdown` before reporting. The implementation and its unit evidence are recorded in the implementation plan §7.5. This does not establish packaged-WebView host acceptance or a real `sd-server` workflow cancellation run, and it does not reopen D-7 or D-8. The frozen order now continues with engine unload/reconfigure during active stages, shutdown joining workflows, then video/workflow overlap including server invalidation and admission transitions; engine lifecycle precedes shutdown integration so a shutdown failure is not attributed to the wrong owner.

`TestCrossWorkloadLeaseSequence` must assert ownership and capacity transitions, not merely admission results. A successful admission does not prove two workloads cannot interfere with the same engine or process. Reclamation refused by an active workflow parent lease produces a structured diagnostic with the refusal reason and workflow engine/claim state; console text alone is non-conforming. The detailed baseline is recorded in `HRA-P1-IMPLEMENTATION-PLAN-2026-11-25.md` §7.4; it does not reopen D-7 or D-8 or authorize a general resource-manager rewrite.

### 12.3 Audio integration gate (required)

Audio is a first-class workload and a required HRA-P1 completion gate. HRA-P1 is not complete while audio admission, lifecycle, cross-workload behavior, or host acceptance remain open. The current `TestASRStartConsultsAdmissionWhileChatLoaded` failure is a policy/test mismatch and stays visible until a documented reserve-contract decision resolves it; it must not be weakened, skipped, or converted into an audio-only admission exception.

The audio work package must establish, through the same authority and lifecycle model as chat, image, and video:

- admission through `resource.Authority`, with any reserve use explained by the documented policy rather than a workload label;
- one owner and one mutation route for ASR/TTS registration, load, execution, cancellation, process termination, and cleanup;
- exactly-once release under completion and cancellation races, with every refusal surfaced;
- accurate pending, resident, and executing resource reporting;
- cross-workload evidence that audio coexists with chat, image, and video without bypassing admission or leaving a process/registration behind;
- actual ASR and TTS host acceptance in packaged WebView2 and the browser dev host, with isolated roots and backend state read back.

The ordered plan is in `HRA-P1-IMPLEMENTATION-PLAN-2026-11-25.md` §7.6. This section authorizes planning and characterization only; it does not change the ASR test, admission policy, or audio production behavior.

The 2026-11-25 ownership audit and its characterization tests are in §7.6.1. They establish that TTS consults the shared authority before a provider is constructed and that a cancelled audio operation reaches the provider call before `Close` reports completion. They also name the open integration gaps: TTS child processes spawned outside job ownership (F-A1), an externally owned MeloTTS service charged to admission (F-A2), registrations that outlive their admission lease without execution revalidation (F-A3), uncancellable Kokoro voice enumeration under a lease (F-A4), and bounded-but-non-forcing audio shutdown (F-A5). Audio remains an open gate; the audit is not acceptance evidence.

F-A1 through F-A5 are now implemented inside the existing architecture (implementation plan §7.6.2). Piper and Kokoro children are job-owned, registered and settled through one waiter; MeloTTS is a Howl-managed loopback service with reference-counted ownership; provider residency is a registration held until a controlled release, with a separate execution reference; enumeration runs as a cancellable managed operation; shutdown escalates to confirmed exit and reports a failure instead of claiming success. The ASR baseline failure is still preserved. A host run and human listening acceptance are now recorded: the operator heard local Piper in their running app and confirmed both read-aloud paths, so voice quality is accepted (`HRA-P1-TTS-LISTENING-ACCEPTANCE-2026-11-25.md` §5). What remains for Gate 8 is integration and lifecycle coverage — cross-workload resource/lifecycle testing, packaged-host audio integration evidence, Kokoro/MeloTTS coverage, and the ASR admission policy/test discrepancy. Gate 8 therefore remains open.

---

## 13. Migration sequence

### Phase 1 — Lifecycle contract and characterization (complete)

- Characterization tests in §11 are written: passing tests document current semantics; known failures are isolated behind `knownfailures`.
- The one-owner, one-process, one-execution-scope, and three-release-event invariants are established as the migration target.
- Every registration creator/remover and every child process is inventoried in §2.
- No lease kind or process registry was added during characterization.
- The transition to production ownership changes is recorded in `HRA-P1-IMPLEMENTATION-PLAN-2026-11-25.md`.

### Phase 2 — Instrumentation only

- Add optional CPU/GPU/temperature/per-process observations.
- Add telemetry age, completeness, and confidence.
- Add process exit observations and stage-boundary observations.
- Do not change any admission decision.

### Phase 3 — Ownership enforcement for all workloads

- Introduce the process registry and preserve existing llama behavior.
- Bring image `sd-server` and video `sd-cli` under it.
- Add load, resident, execution, and workflow leases, including ASR/TTS load and execution.
- Route registration creation/removal through the authority.
- Replace direct `Engine.Shutdown` unregistration and `stopServer` direct `Wait`.
- Add workflow parent scope and component claims.
- Keep audio on the same reserve policy as every other workload; no audio-only admission branch.

### Phase 4 — Unified incremental accounting

- Use `ProjectedPeakVRAMMB` for the load commitment.
- Add incremental execution commitments.
- Keep pending until reflected or released.
- Remove the separate `resource.Allocate` policy.
- Include RAM-only blocking.
- Use one reserve policy.

### Phase 5 — Engine adapters and video completion

- Video forwards managed budget and hardware flags.
- Video workload includes frames, fps, resolution, VAE, and temporal tiling.
- Adapters report placement, peak, and exit evidence.
- Image/chat/audio adapters report actual process and placement evidence.

### Phase 6 — Persistent observation and calibration

- Persist observations only after phases 1–5 pass.
- Calibration may improve prediction; it never grants permission.
- Reducing safety margins requires stronger evidence than raising them.

---

## 14. Verification basis and limitations

### 14.1 Baseline test results

The following were run with `-count=1` and passed:

```text
go test -count=1 ./internal/resource ./internal/host ./internal/image/engine/sdcpp
ok  howl-chat/internal/resource
ok  howl-chat/internal/host
ok  howl-chat/internal/image/engine/sdcpp

go test -count=1 . -run 'TestImageAdmission|TestChatAdmission|TestFailedUnload|TestHandoff|TestStaleTelemetry|TestLoadedChatRefusesImagePreflight|TestReverseHandoff'
ok  howl-chat
```

Gate 1 focused evidence:

```text
go test -count=1 ./internal/resource
ok  howl-chat/internal/resource
```

These establish existing behavior only. They do not prove the new contract, the races in §2, or the cross-workload scenario.

### 14.2 Verified limitations

- Passing characterization tests now exist in `internal/resource/characterization_test.go`.
- Gate 1 graduated D-1 and D-9 from `knownfailures` into the default resource suite; five known failures remain isolated in `internal/image/engine/sdcpp` and the root `main` package.
- `TestPendingCommitmentSurvivesUntilObservation`, `TestVideoExecutionHasLease`, `TestVideoLeavesNoImageRegistration`, `TestWorkflowReusesEngineAcrossStages`, `TestWorkflowRetainedUpscalerIsClaimed`, `TestOptionalTelemetryAbsenceDoesNotRefuse`, `TestManualBudgetProvenanceSurvivesReload`, and `TestPublicRegistrationBypassesAuthority` are still planned.
- No phase-one acceptance tests in §12 exist yet.
- `internal/audio`'s pre-existing `TestASRStartConsultsAdmissionWhileChatLoaded` remains red: its 4500 MB snapshot admits the current 4096 MB ASR profile through the soft-reserve branch. Gate 1 did not change that branch. The baseline and policy classification are recorded separately in `HRA-P1-AUDIO-ADMISSION-BASELINE-2026-11-25.md`.
- No user-visible workflow was run for this contract.
- Gate 1 changed `internal/resource/admission.go`; no other production package has started.
- No persistent calibration was performed.

### 14.3 Defect register

| ID | Defect | Evidence |
|---|---|---|
| D-1 | Grant under-reserves: admission checks `ProjectedPeakVRAMMB`, `grant` stores `RequiredVRAMMB` | `internal/resource/admission.go` `grant` / `needVRAM` |
| D-2 | Pending reservation releases before physical telemetry confirms the allocation | `registerImageModel`, `registerAudio`, chat register/release |
| D-3 | Chat monitor can register after `UnloadModel` | `app-model.go` health-success path |
| D-4 | Failed image reload can leave the previous registration | `ensureEngineReadyLocked` error path |
| D-5 | Image shutdown/unload can race generation callbacks; jobs are not joined | `app-lifecycle.go`, `image_reclaim`, `Engine.Shutdown` |
| D-6 | Video `stopServer` calls `Wait()` despite the single-reaper invariant | `video.go:stopServer`, `supervisor.go:startProcessReaper` |
| D-7 | Video has no resource registration/execution lease and can leave an image registration behind | video branch, `reclaimImageEngineAfterGeneration` early return |
| D-8 | Workflow stage 4 assumes the upscaler is retained without a claim | `app-image-workflow-run.go` stage 4 |
| D-9 | RAM-only residents are omitted from `blockingModels` | `internal/resource/admission.go` |

**Gate 1 remediation.** D-1 and D-9 are closed for the accounting and blocker slices. `grant` commits `committedVRAMMB()`, a resident stage's projected peak includes its newly added adapter bytes, and `blockingModels` includes RAM-only residents when the request has a RAM requirement. The original assertions pass in the default suite, and `TestCommitmentUsesComparableProjectedPeakAcrossPaths` pins full-load, resident-stage, resident-stage-with-adapter, and replacement commitments against `Required + Projected` double-counting. Replacement credit observation remains a later `IP-8` requirement.

**D-7 remediation.** Video holds its admission lease across `SelectModel` and `GenerateVideo` and releases it exactly once after the CLI path returns. The `sd-server` stop before a video run reports process exit through the app hook, which removes the stale image registration. `sd-cli` is adopted into the process registry immediately after start and released after its recorded waiter observes exit, including cancellation. This is the named D-7 slice, not a full authority registration transaction or a complete video shutdown matrix.

**D-8 remediation.** A workflow now owns one `resource.WorkflowLease` parent scope, with a stable engine-instance identity. Stage 2 acquires the upscaler component claim; stage 4 verifies the same key and refuses when it is absent. The key includes engine instance and process identity, role, normalized path, weight type, and companion set, and claims are reference-counted and idempotent on release. Reclamation is refused while the parent scope is active, preserving the retained component across batch item boundaries, and a replaced engine instance invalidates the scope instead of being used silently. This is lifecycle ownership, not a new admission-capacity lease; the complete user-visible workflow remains Gate 7 evidence.

---

## 15. Phase-one acceptance boundary

Phase one is complete only when all four are demonstrated:

| Property | Evidence |
|---|---|
| Ownership | Every registration and child process has one owner; no registration outlives its process |
| Accounting | Incremental peak is reserved; pending/reflected/released are distinct; no double charge |
| Enforcement | Concurrent requests cannot spend the last capacity; shared components are claimed; no path bypasses the authority |
| Recovery | Crash, cancel, shutdown, and process exit leave no orphan process, phantom registration, or stale reservation |

A successful operation with an open lease is not a pass. A passing unit suite without the cross-workload scenario is not phase-one completion.

### 15.1 User-visible gate

When phase one changes user-visible behavior, the accepting workflow must be exercised in both supported hosts, with an isolated state root:

- packaged desktop app (`desktop,production`) with real mouse/screenshots;
- browser dev host (`uiharness`) over the real backend.

Native-surface differences must be recorded explicitly. The workflow must demonstrate: chat resident, image generation, video admission/execution, image completion, and clean shutdown without a phantom model or orphan process. Until that is run, phase one is not user-visible complete.

---

## 16. Implementation planning record

The production ownership changes, defect reconciliation, invariant coverage, and acceptance evidence for §11 are recorded in:

- `docs/architecture/HRA-P1-IMPLEMENTATION-PLAN-2026-11-25.md`

That plan does not change this contract's acceptance boundary. It maps each known failure and untested invariant to an owning package and requires default-suite evidence before any known-failure tag is removed. No production instrumentation was added to make the remaining tests possible.
