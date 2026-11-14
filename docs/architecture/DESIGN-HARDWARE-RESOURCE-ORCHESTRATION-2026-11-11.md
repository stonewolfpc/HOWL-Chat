# Design — hardware detection and resource orchestration

**Status:** Proposed architecture design; no implementation authorization  
**Date:** 2026-11-11  
**Owner input:** `ODA-09` (automatic hardware detection and optimization are core), `UD-01` (per-chat models)  
**Plan sources:** `REC-M-03`, `REC-M-07`, `REC-M-08`, `REC-M-15`, `REC-M-16`, `CWP-15a/15b`, `CWP-08a/08b`  
**Related design:** `DESIGN-MEMORY-CONTINUITY-2026-11-11.md`, `RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`

## 1. Decision in one sentence

Howl will own a vendor-neutral hardware inventory and a single resource authority: it detects what it can, reports confidence honestly, reserves at least 10% of currently free VRAM by default, admits or refuses loads through one lease model, and allows an explicit, warned override rather than pretending an unknown system is safe.

## 2. Current implementation and its limits

| Surface | Current behavior | Limit |
|---|---|---|
| Snapshot | `internal/resource` exposes GPU/RAM totals and free memory | No CPU model, storage headroom, driver, backend, or per-device detail |
| NVIDIA | `nvidia-smi` query, then registry estimate, then WMI | `nvidia-smi` is an external process; WDDM reporting differs from TCC; fallback can be stale |
| AMD / Intel | Registry/WMI discovery and name/size estimates | No direct vendor telemetry; live used VRAM is weak or absent |
| Safety margin | Allocator uses 90% of reported free VRAM | The 10% reserve is implicit, not a named product rule |
| Context/KV | GGUF profile plus `context_size` and `gpu_layers` estimator | No general KV-cache model, batch, cache type, or offload matrix |
| Consumers | Chat and image arbitration exist; TTS is outside admission; image engine is not consistently registered | Multiple authorities and unregistered allocations |
| Override | `force_load` can bypass advisory checks | No forced-override audit path or recovery contract |
| Image estimates | Fixed model-size and overhead guesses | Estimates are not reconciled with live allocation |
| Release | Fixed 1.5 s wait after unload | A sleep is not evidence that VRAM was reclaimed |
| Failure | Fallback to last-known values; low-confidence advice | Does not distinguish “unknown” from “free”; does not define user override behavior |

The existing manager is a useful **advisory foundation**. It is not yet the single resource authority required by `REC-M-07`.

## 3. Design principles

1. **Detect, do not assume.** Hardware facts carry source, timestamp, and confidence.
2. **Unknown is not zero.** A failed telemetry query must not become “no GPU” or “free VRAM.”
3. **One admission authority.** Every model consumer obtains a lease from the same authority.
4. **Reserve before allocate.** At least 10% of currently free VRAM stays unused by default.
5. **Fit is a prediction.** The UI may say “expected to fit,” never “will fit,” without evidence.
6. **Degrade conservatively, not silently.** If monitoring is unavailable, warn and allow an attempt rather than blocking the user.
7. **Let the user override the prediction, not the record.** A forced attempt is logged and visibly marked as an override.
8. **Model identity is not a filename.** Architecture, quantization, template, backend, and hardware compatibility are separate facts.

## 4. Hardware inventory

```text
HardwareInventory
  captured_at
  os / arch
  cpu:
    model, physical_cores, logical_cores, ISA features, numa?
  memory:
    total, available, committed?
  storage:
    model_roots[], free_bytes, read_only?
  gpus[]:
    stable_device_id, vendor, name, driver, backend
    total_vram, used_vram, free_vram, reserved_vram?
    memory_bus?, unified_memory?
    telemetry_source, confidence, last_error
  monitors:
    provider_status, permissions, missing_fields
```

Per-GPU allocation is tracked separately. Multi-GPU is a set of devices, not a summed pool, because layer placement and peer access are not interchangeable.

## 5. Telemetry provider abstraction

```text
GPUTelemetryProvider
  Probe() -> ProviderStatus
  Snapshot(device_id) -> GPUMemoryReading
  Capabilities() -> {total, used, free, processes, temperature, power, utilization}
```

### 5.1 Provider order

| Priority | Provider | Use | Failure behavior |
|---:|---|---|---|
| 1 | NVIDIA NVML | Total/used/free, process attribution where available | Dynamic load failure falls through |
| 2 | AMD ADL/PDH | Adapter identity and dedicated usage; OS GPU counters by LUID where available | Partial result is marked partial, not complete |
| 3 | Intel Level Zero Sysman / PDH | Device memory where exposed | Missing counters are explicit |
| 4 | DXGI/WMI/registry | Discovery and static capacity | Never treated as live used/free unless the counter is known |
| 5 | Process accounting | Sum of Howl-owned allocations | Used only for Howl's own processes, not total system use |
| 6 | Last-known snapshot | Short-lived continuity after a transient failure | Stale; lower confidence; visible |

Vendor libraries are optional and dynamically loaded. No vendor SDK is required to start Howl. A provider that returns partial data is preferred over a provider that fabricates a complete picture.

**Windows caveat:** NVML documents that WDDM memory reporting differs from Linux/TCC; AMD's own VRAM usage API has known unreliable results on some adapters, so OS dedicated-memory counters are a fallback, not a guarantee. The design therefore treats every reading as a claim with confidence.

## 6. The 10% VRAM reservation

For each GPU, at admission time:

```text
free_now = live_free_vram if confidence >= threshold else unknown
reserve = max(0.10 * free_now, minimum_reserve)
usable = max(0, free_now - reserve - estimated_runtime_overhead)
```

Rules:

- The reserve is calculated from **currently free** VRAM, not total VRAM.
- The reserve is recalculated at commit; a fit computed before another process allocated memory is stale.
- A lease cannot spend the same free bytes twice.
- `minimum_reserve` is a small absolute floor for driver/context churn, not a substitute for the 10%.
- Runtime overhead (CUDA context, graph buffers, image encoders/VAE, audio buffers) is estimated separately and shown in the fit report.
- If `free_now` is unknown, Howl does not claim a fit and does not auto-evict.

The existing allocator's 90% constant is compatible with this rule, but the design makes the reserve explicit, per-device, and auditable.

## 7. Model suitability

Suitability is a report, not a boolean.

```text
ModelFitReport
  model_identity
  backend
  quantization
  context_tokens
  batch_tokens
  kv_cache_bytes
  weights_bytes
  runtime_overhead_bytes
  required_vram
  free_vram
  reserve
  usable_vram
  offload_plan
  expected: fits | tight | cpu-fallback | insufficient | unknown
  confidence
  blockers[]
  alternatives[]
  warnings[]
```

### 7.1 KV-cache estimate

The estimator must use model metadata, not a generic token count:

```text
kv_bytes ≈ 2 × layers × kv_heads × head_dim × context_tokens × bytes_per_element × batch
```

plus attention/backend workspace. For grouped-query attention, `kv_heads` is not `attention_heads`. For sliding-window or recurrent models, the formula is replaced by a model-family adapter. If the architecture is unknown, the report says `unknown` and gives a conservative range.

### 7.2 Compatibility state

| State | Meaning |
|---|---|
| `stored` | A setting exists |
| `sent` | It was passed to the backend |
| `supported` | The backend claims support for this model/version |
| `effective` | It was observed in an authorized run |

No design may promote a stored flag to “supported” or “effective.” That is the core correction required by `REC-M-03` and `REC-M-13`.

## 8. One resource authority

The resource authority owns admission, leases, release, and reconciliation.

```text
PlanLoad(request) -> ModelFitReport
AcquireLease(request) -> lease_id | refusal
CommitLoad(lease_id) -> loaded_model
Release(lease_id) -> release_receipt
Reconcile() -> actual_vs_tracked
```

### 8.1 Consumers

All of these must request a lease:

- chat/LLM;
- image/diffusion;
- TTS/ASR;
- embedding/index (only if RAG is enabled under the relevant authorization);
- future studio modalities.

TTS/ASR is a **resource consumer, not a product studio**. Character voice belongs to Chat (`UD-05`); it still requests the same lease as any other model workload.

A model loaded outside the authority is an **orphan** and is surfaced in diagnostics. The authority does not silently kill an orphan; it reports it and offers a controlled recovery path.

### 8.2 State machine

```text
unloaded → planning → admitted → loading → ready → busy → draining → unloading → unloaded
                                  ↘ failed
```

- `admitted` does not reserve indefinitely; leases expire.
- `loading` can fail without corrupting the previous ready model.
- `busy` cannot be evicted without drain.
- `unloading` waits for the process and driver to reclaim memory.
- `released` requires evidence or an explicit “release unverified” state.

### 8.3 Per-chat model binding

`UD-01` makes the selected model a chat property. The binding is resolved before a turn:

1. read the chat's model binding;
2. check model availability and compatibility;
3. request an admission lease;
4. if the model is not resident, load at a turn boundary or offer a safe switch;
5. if the lease fails, show the fit report and alternatives;
6. never silently substitute a different model.

Hot-swapping during an active stream is convenience, not a release blocker. A safe boundary is required; a forced stop is visible and recoverable.

### 8.4 Atomic lease and commit

Admission is a two-phase reservation, not a check followed by an unrelated load.

```text
Lease
  lease_id
  idempotency_key
  owner: chat | image | tts | asr | embedding | other
  chat_id / job_id
  device_ids[]
  reserved_bytes
  priority
  state: reserved | committed | draining | released | expired | failed
  expires_at
  observed_generation
```

1. **AcquireLease** validates the request, reads one coherent device snapshot, computes the fit including the 10% reserve, and atomically records the reservation under the authority's write lock.
2. A reservation is visible to every other consumer immediately. The same free bytes cannot be promised twice.
3. **CommitLoad** rechecks that the snapshot and reservation are still valid. A material change rejects or re-plans the lease; it does not load against stale numbers.
4. A load that fails after allocation transitions to `failed`, reconciles actual usage, and releases only what is proven to be owned by that lease.
5. Retries reuse the idempotency key. A duplicate request returns the existing reservation or receipt rather than creating a second allocation.
6. Reservations have a bounded lifetime and a renewal/heartbeat path. Expired reservations are reclaimed, not left as phantom capacity.
7. A crash leaves reservations reconcilable from durable lease records plus a fresh device snapshot.

### 8.5 Concurrency, contention, and busy models

| Situation | Required behavior |
|---|---|
| Two requests fit independently but not together | First committed reservation wins; the second is refused or re-planned with a visible fit report |
| Requests target the same chat/model | Coalesce compatible work; serialize changes at a turn boundary |
| Requests target different chats on one device | One authority arbitrates; no per-chat allocator |
| Active generation is `busy` | Do not auto-evict. Queue within policy or refuse with an explanation; a user-requested forced stop is visible and recoverable |
| A higher-priority request arrives | Apply a declared priority policy with fairness/aging; never starve a waiting consumer indefinitely |
| Android requests capacity | They use the same lease path and desktop policy; remote origin grants no priority |
| A lease expires while queued | Release the reservation; the request may retry with fresh fit data |
| Two consumers race at commit | The revision/commit check admits one; the loser re-plans or fails explicitly |

`busy` is not a synonym for “safe to kill.” A drain boundary is consumer-specific: stop accepting new work, finish or cancel the active stream under the declared policy, then unload. Chat, image, TTS/ASR, embedding, and future studios must each define that boundary.

### 8.6 Release verification and rollback

- Release proceeds through `draining → unloading → released`; no fixed sleep is accepted as proof.
- The authority polls or reconciles the process, driver-visible usage, and its own ledger over a bounded window.
- If usage is still above the expected baseline, report **release unverified** and retain the reservation/ledger entry.
- If the prior model was resident and remains safe, keep it; otherwise report the unloaded state rather than claiming a model is ready.
- Failed prefetch/load must not destroy a working prior model unless the user explicitly chose the switch.
- Orphan processes and unregistered allocations are reported; recovery is controlled and attributable.
- A forced override records the warning, the observed outcome, and whether release was verified; it never becomes a silent default.

### 8.7 Live changes and existing allocations

- Context size, KV-cache type, batch size, quantization, offload, backend, or device selection changes invalidate the fit. Re-plan and, where needed, unload/reload; never patch admission in place.
- Existing non-Howl GPU allocations are part of the observed free-memory baseline, not free capacity. Howl must not assume its registry knows them.
- On platforms where live usage cannot be attributed per process, say so; do not convert an unexplained difference into an optimistic budget.
- The 10% reserve is applied to currently free VRAM after existing allocations, not to a historical or theoretical total.
- Reconcile orphaned allocations before granting a new large lease; a stale ledger must not outvote the device.

## 9. Fallback and override

### 9.1 Monitoring unavailable

When live VRAM is unavailable:

- show `telemetry: unavailable` and the reason;
- use last-known total/used with a stale timestamp, or process-accounted usage;
- label the fit `unknown`, not `fits`;
- allow the user to attempt the load;
- do not auto-unload another model;
- keep the 10% rule when a free value is available; when no free value exists, do not fabricate one;
- record the attempt and outcome for future calibration.

### 9.2 Conservative degradation

- Prefer partial offload or smaller context when it preserves usability.
- Do not silently reduce output quality without showing what changed.
- Use CPU fallback when allowed, but show expected speed/memory consequences.
- Refuse only when the user has not allowed the fallback or the backend cannot operate.

### 9.3 Power-user override

The normal path is automatic. Overrides are explicit:

- require a confirmation dialog with the exact predicted shortfall;
- show which safety margin is being overridden;
- preserve the 10% reserve as the default and mark any bypass;
- record a local override event;
- do not turn the override into a permanent default;
- never claim the load was safe or verified.

## 10. UI contract

The user sees a **readiness card** before a load:

- model name, quant, context, backend;
- expected fit and confidence;
- GPU and VRAM summary, with telemetry source and age;
- KV/weights/overhead breakdown;
- what will be unloaded, if anything;
- alternatives: smaller context, fewer layers, CPU fallback, another model;
- “attempt anyway” with warning when monitoring is unavailable or fit is tight;
- explicit “unknown” when data is missing.

No raw “VRAM free: 0 MB” is shown merely because a query failed.

## 11. Failure and recovery

| Failure | Behavior |
|---|---|
| Vendor library missing | Use next provider; partial/unknown data, not fatal |
| Telemetry timeout | Keep last-known values with age; mark low confidence |
| Driver reports impossible total | Reject and retain last known good; show error |
| Load fails after allocation | Roll back to the prior safe model when possible; reconcile actual usage; never leave a phantom ready state |
| Unload does not reclaim | Surface “release unverified”; retain the ledger entry; do not claim freed |
| Orphan model | Report; require controlled recovery |
| Two leases exceed free VRAM | Atomic reservation admits one; reject or re-plan the other; never double-spend |
| Lease expires during a load | Cancel/reconcile; do not commit against an expired reservation |
| Context/KV setting changes | Invalidate fit; re-plan or restart; no in-place admission |
| Existing non-Howl allocation appears | Treat it as used; re-plan rather than assuming the registry is complete |
| User forces load | Warn, record override, observe outcome; no safety claim |

## 12. Interface and ownership boundaries

| Component | Owns | Does not own |
|---|---|---|
| Hardware inventory | Detection and confidence | Model selection |
| Telemetry providers | Vendor-specific readings | Admission |
| Resource authority | Leases, admission, release | Model quality claims |
| Model manager | Process lifecycle | Global policy |
| Chat | Per-chat binding | Hardware totals |
| UI | Explanation and override | Silent substitution |
| Memory/RAG | Their own workload requests | Resource policy |

## 13. Acceptance design — not run

1. Provider absence produces `unknown`, not zero.
2. A 10% reserve is visible and cannot be spent by default.
3. Two concurrent leases cannot double-spend free VRAM.
4. A stale fit is rejected or re-planned at commit.
5. Unload reports verified versus unverified release.
6. Every consumer appears as a lease or an orphan.
7. Per-chat model binding resolves correctly after restart.
8. Low-end/CPU fallback is explained before load.
9. Unknown telemetry still allows a warned attempt.
10. A forced override is logged and does not become a silent default.
11. KV estimates are checked against model metadata, not only filename.
12. No benchmark claim is made without named hardware and an authorized run.
13. Concurrent acquire/commit races cannot reserve the same bytes twice.
14. An expired or crashed lease is reclaimed without freeing another lease's bytes.
15. A failed load restores the prior safe model or reports an explicit unloaded state.
16. Release is verified by observation, not by a fixed sleep.
17. Context/KV/offload changes force a new fit decision.
18. Every consumer, including TTS/ASR and Android-requested work, is a lease or a reported orphan.
19. A busy model is drained or refused according to policy, not silently killed.

No telemetry test, model load, benchmark, or hardware experiment was performed.

## 14. Open decisions

| ID | Decision | Why open |
|---|---|---|
| `DA-H1` | Provider implementation order and which vendor libraries to link dynamically | Runtime/package constraints; no experiments authorized |
| `DA-H2` | Exact `minimum_reserve` and overhead model | Needs measurement, but no benchmark authorized |
| `DA-H3` | Whether multi-GPU is layer-split, tensor-parallel, or unsupported per backend | Backend-dependent |
| `DA-H4` | Model compatibility matrix source and maintenance | Requires evidence, not filename heuristics |
| `DA-H5` | Default low-end policy: CPU fallback versus refusal | Product experience tradeoff |

## 15. External references

- NVIDIA NVML API: https://developer.nvidia.com/management-library-nvml
- NVML memory/WDDM note: https://docs.nvidia.com/deploy/nvml-api/latest/api/group__nvmlDeviceQueries.html
- AMD ADL/PDH and Windows telemetry discussion: https://github.com/nalilord/AMDPlugin
- Windows GPU memory counters via PDH/LUID (practitioner evidence): https://www.stratat.com/ja/pulls/1434

These references inform the design; they do not establish Howl runtime behavior.

## 16. Non-goals

- Overclocking, fan control, voltage, or persistence-mode changes.
- Guaranteeing a model will fit.
- Treating a vendor SDK as always installed.
- Auto-evicting a busy model.
- Claiming release without observation.
- Running any hardware or benchmark investigation under this design.
