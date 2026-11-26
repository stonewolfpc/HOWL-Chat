# Milestone 2 — ownership audit (authority and lifecycle)

Date: 2026-11-26
Status: audit only. No production code changed, nothing deleted.
Baseline: `f131f09` (Milestone 1 evidence)

## Purpose and method

Milestone 2 is treated as an ownership audit before a cleanup exercise. The goal
is one authoritative answer to five questions, not fewer lines of code:

1. What resources are physically available?
2. What capacity is already reserved?
3. What workload owns each allocation?
4. Who may admit, reclaim, or release it?
5. When is reclaimed capacity actually safe to reuse?

Method: read every owner and every caller in `internal/resource`, the four load
paths (`app-model.go`, `app-image-model.go`, `internal/audio/pipeline.go`,
`app-image-workflow-run.go`), the two reclaim paths, and the bindings in
`app-resource.go` / `app-vram-arbitration.go`. Findings are from code and call
graph, not from runtime experiments; the one behavioural claim that would need a
reproduction (the cross-workload interleaving in §5) is marked as unverified.

## 1. What resources are physically available?

Owner: `resource.ResourceManager` (`internal/resource/manager.go`), via
`RefreshSnapshot()`, which merges two independent detectors:

| Value | Detector | Notes |
|---|---|---|
| live GPU total/used/free per device | `internal/resource/gpu_info_windows.go` (`queryGPUs`) | sets `GPUInfo.Live`, `HardwareSnapshot.VRAMLive` |
| GPU total VRAM, GPU presence, RAM, tier | `internal/audio/hw/gpu_windows.go` → `hw.DetectHardware()` | nvidia-smi `memory.total`, registry, WMI |
| RAM total/available | `hw.DetectHardware()` | |

Admission trusts only `VRAMLive`; when live used-memory is unavailable it
refuses (or takes the manual chat rung) instead of reusing a stale value —
`admission.go`, `if totalVRAM > 0 && !snap.VRAMLive`.

Competing owners found: total VRAM has two readers (resource detector and
`audio/hw`), with `ResourceManager` as the single merge point. Not yet a defect,
because only the merge point is consulted for decisions — but it means a future
caller that imports `audio/hw` directly can answer question 1 differently. The
`last-known-good` fallback in `refreshHardware` also rewrites `TotalVRAMBytes`
while clearing `VRAMLive`, which is correct for admission but makes the snapshot
self-inconsistent for display consumers (`app-resource.go` reports it plainly).

## 2. What capacity is already reserved?

There is no single number. Four reserve computations exist:

| # | Where | Rule | Consumer |
|---|---|---|---|
| R1 | `admission.go` `DefaultReservePercent` | 10% of **effective free** VRAM, RAM; plus a 500 MB soft allowance above a 128 MB floor | every admission decision |
| R2 | `allocator.go` | 90% of free VRAM / 95% of free RAM (`vramSafetyPct`, `ramSafetyPct`) | `resource.Allocate` → `RequestAllocation` binding only, no UI caller |
| R3 | `app-image-reclaim.go`, `app-image-placement.go` | 10% of **device total** | image recycle trigger and placement report |
| R4 | `app-image-model.go` / `imageVRAMBudgetMB` | the user's `--max-vram` cap, and for images a device-total 10% floor substitution inside admission | sd-server startup and the image reclaimer floor |

R1 and R3 differ in base (free versus total) and exist because each was tuned
against a different failure. They are not interchangeable, and neither is
provably right until capacity inputs are trustworthy — which is why this
milestone does not reconcile them.

## 3. What workload owns each allocation?

Declared owner: `ResourceManager.models`, keyed `type:name`.

Registered by:

| Workload | Route | Footprint source |
|---|---|---|
| chat | `registerChatModelAfterHealth` (`app-model.go`) | `loadDecision.Estimate.RequiredVRAMMB` |
| image | `registerImageModel` (`app-image-model.go`) | `decision.Estimate.RequiredVRAMMB` |
| audio | `internal/audio/pipeline.go` after `Admit` | the request's `AudioVRAMMB`/`AudioRAMMB` |

Competing owners found:

- **Engine truth versus registry.** Each engine also knows what it is running:
  `llama` state via `IsModelLoaded`, sd-server via `sdcpp.Engine.GetLoadedModel`
  and `images.LoadedOptions()`. `loadedImageModelInfo` (`app-vram-arbitration.go`)
  reads the *name* from the engine and the *VRAM* from the registry — one answer
  composed from two owners. The comment there already concedes that the image
  engine is not tracked by the registry, yet the registry is still consulted for
  the number.
- **Estimate registered as fact.** Both chat and image register
  `RequiredVRAMMB`, not the committed figure. For chat those are equal (no
  projection), but for image `committedVRAMMB()` is `ProjectedPeakVRAMMB`
  (resident + activation + overhead). The ledger therefore understates a
  resident image engine by its activation and overhead — exactly the gap the
  HiDream-class case sits in.
- **A frontend write path into the ledger.** `RegisterLoadedModel` and
  `UnregisterLoadedModel` (`app-resource.go`) are Wails bindings that write the
  resident ledger directly. No current frontend code calls them, but they are a
  second mutation route into the single authority's input, available to any page
  script.

## 4. Who may admit, reclaim, or release?

Admission: one owner, `resource.Authority.Admit`, with an in-flight ledger
(`inFlight`) and leases. Verified ordering on every load path:

| Path | Admit | Register | Lease released | Source |
|---|---|---|---|---|
| chat | `app-model.go:303` | after health check | after registration (`:593`) | correct order, no double-count window |
| image (generation) | `app-image-model.go:256` | after `EnsureReady` | after registration (`:145`) | correct |
| image (video execution) | `app-image-model.go` | after `SelectModel` | lease handed to the caller, released after sd-cli exits | correct |
| audio | `pipeline.go:765` | `:783` | `:792` | correct |
| preflight (`CheckModelFit`, `CheckImageFit`, q8 probe) | yes | no | released immediately | correct, read-only |

Reclaim has three separate routes:

| Route | Trigger | Wait rule | Re-admits? |
|---|---|---|---|
| chat takes the device (`FreeVRAMFor` → `reclaimOther` → `waitForAdmissionFit`) | chat load wants GPU while image is resident | polls `Admit` against fresh snapshots until admitted or 10 s | yes, same immutable request |
| image takes the device (`FreeVRAMForImage`) | generation preflight, user-approved | same `reclaimOther` | yes |
| image recycles itself (`reclaimImageEngineAfterGeneration`) | post-generation measurement breach or abandoned state | none; reloads `loaded` options immediately | **no** |

Release of tracked state happens through `UnregisterModel` from the unload
paths (`UnloadModel`, `UnloadImageModel`, `invalidateImageRegistrationsIfStopped`,
`failChatLoad`), and through `UnregisterLoadedModel` from the frontend binding.

## 5. When is reclaimed capacity actually safe to reuse?

Today: not guaranteed. This is the milestone's most consequential finding.

- There are two independent start locks: `modelStartMu` (chat, `app.go:51`) and
  `imageStartMu` (image, `app.go:54`). They never nest or coordinate.
- `UnloadModel` (`app-model.go:731`) — the call `FreeVRAMForImage` makes to free
  the device for image work — takes no lock. It can stop llama-server while a
  chat load holds `modelStartMu` and is between spawn and health check.
- `UnloadImageModel` (`app-image.go:150`) — the call chat's handoff makes — takes
  no lock. It can shut down sd-server while an image workflow holds
  `imageStartMu` and is mid-stage, or while the reclaimer is reloading.
- The image reclaimer reloads weights with `engine.LoadModel(loaded)` without a
  fresh admission decision. Its comment says the options were already admitted,
  but that admission happened before an arbitrary amount of intervening work; if
  another workload took VRAM in between, the reload either fails or is re-fitted
  by sd.cpp and lands somewhere the plan never described.
- The one place that does answer the question correctly is `waitForAdmissionFit`:
  it never treats a successful kill as proof of release, and it re-admits the
  same request against live telemetry. That is the shape the other routes need.

Not reproduced at runtime; the claim is from the call graph and the absence of
any lock acquisition on those two unload paths. A deterministic test with fake
engines can settle it without a GPU.

## Built but unreachable (inventory, not a deletion list)

| Symbol | Callers | Why it matters |
|---|---|---|
| `gguf.ValidateSpawn` / `ValidateSpawnLegacy` (`memory.go:280`) | none, including tests | duplicate spawn-safety policy; would disagree with admission if revived |
| `resource.Allocate` (`allocator.go`) | only `App.RequestAllocation` | second reserve policy (R2) and a second eviction ladder |
| `App.RequestAllocation` (`app-resource.go:123`) | Wails binding, no UI caller | exposes R2 as if it were an authority |
| `App.RegisterLoadedModel` / `UnregisterLoadedModel` | Wails bindings, no UI caller | writable ledger, see §3 |
| `AdmissionDecision.Warning` on the manual rung | consumed by `app-model.go` | live, not unreachable — listed to avoid conflating the two |

Per the agreed sequence these are **not** removed in this milestone. Their
callers are now enumerated; the decision to delete follows the reconciliation,
not the other way round.

## What Milestone 2 should unify, in dependency order

Each step names the evidence that would justify it, so none can be closed on
argument alone.

1. **One lifecycle gate for capacity hand-offs.** A single cross-workload
   acquire/release scope that `UnloadModel` and `UnloadImageModel` must hold
   before stopping another workload's process — or, equivalently, make every
   stop go through `reclaimOther`'s confirm/wait shape. Evidence: a
   deterministic interleaving test (chat spawning while image frees the device,
   and the reverse) plus a two-host run of chat↔image handoff with both logs
   showing one process stopped only after its owner released it.
2. **One resident ledger fed by engine truth.** Register the committed footprint
   (`committedVRAMMB`) and, where the engine reports it, the observed buffer
   total; remove the registry/engine split in `loadedImageModelInfo`. Evidence:
   the observed image resident total from a real generation compared with the
   registered value — that measurement is Milestone 3's first deliverable, which
   is why step 2 cannot complete before it.
3. **One reclaim safety rule.** Re-admit before any reload (the image reclaimer's
   recycle), or justify the exception with a measurement showing the reload
   cannot exceed the earlier commitment. Evidence: a recycle under a competing
   workload.
4. **One reserve.** Only after (2) and (3) produce measured residents: reconcile
   R1/R3 into one policy with an explicit base (free or total) and one soft
   allowance owner, then decide whether R2's safety percentages are still a
   different policy or a stale copy.
5. **Retire duplicates.** Delete or re-scope `ValidateSpawn`, `Allocate`,
   `RequestAllocation`, and the ledger-writing bindings once every caller above
   is accounted for and step 1–4 have landed.

Steps 1 and 3 are independent of measurement and can proceed now. Step 2 is
partly blocked on Milestone 3. Step 4 is fully blocked on it. This is the
dependency surfacing naturally rather than being imposed.

## Verified limitations of this audit

- No runtime experiment was performed; §5's interleaving is a code-graph claim.
- `internal/audio` was read at its admission/registration seams only; its
  internal residency policy (`tts_residency.go`) was not audited in depth.
- The Windows and Unix GPU detectors were read for call shape, not behaviour.
- No UI surface was exercised; this document makes no claim about visible
  behaviour.
