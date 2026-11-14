# Amendment 01 — unified hardware intelligence, model residency, and workload orchestration

**Date:** 2026-11-11  
**Status:** Planning amendment; no architecture approval, implementation, migration, runtime, model, benchmark, or network authorization  
**Amends:** `docs/architecture/INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md`  
**Original plan SHA-256:** `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf`  
**Original plan size:** 53,197 B / 724 lines  
**Preservation rule:** The original plan and its `.sha256` sidecar remain byte-identical. This amendment is the controlling delta.

## 1. Amendment decision

The previous plan treated hardware and resource orchestration as one important research workstream among several. The owner has clarified that it is a **core, system-wide product architecture**:

> The user selects what they want to use and requests the work. Howl determines how to execute it safely and efficiently.

This is not merely a model loader, a VRAM warning, or a set of independent modality-specific loaders. It is a unified authority over:

- hardware discovery and telemetry;
- existing allocations, including allocations outside Howl;
- model suitability and placement;
- model selection versus residency;
- workload queuing and admission;
- load, readiness, execution, drain, and unload transitions;
- resource release verification;
- scheduling across concurrent work;
- failure, uncertainty, and manual override behavior.

The amendment preserves the product direction already accepted:

- **Chat** is the first stage, including character voice/TTS and Narration.
- **Image** is next.
- **Music and 3D** remain future dedicated studios.
- TTS remains a Chat capability, not a separate studio.
- Broad UI unification remains later and does not block usability, accessibility, or functional corrections.

## 2. What changed and why

| Plan v1 area | Amendment |
|---|---|
| R-03 was one of six research streams | R-03 is elevated to the central hardware-intelligence program, with sub-workstreams R-03a–R-03j |
| Resource authority described leases and release | The authority now owns placement, residency, scheduling, readiness, transition, and release verification |
| Per-chat model binding was explicit | Model **selection**, **residency**, and **execution** are now separate states; tab switching does not unload |
| Concurrency was a failure-mode concern | Concurrent residency and concurrent execution are separate design problems and explicit acceptance areas |
| Multi-GPU was an open backend question | Multi-GPU is a first-class requirement; per-backend limits must be evidenced, not assumed away |
| CPU/RAM offload was a fallback | CPU/RAM placement is a legitimate strategy chosen from supported capability, performance, resources, and user preference |
| 10% VRAM reserve was a safety rule | The reserve remains authoritative, applied after existing allocations, with additional peak/fragmentation margins still under research |
| Release was verification-aware | The lifecycle now explicitly distinguishes backend unload, observed release, partial release, and unknown release |
| Manual overrides existed as a principle | Overrides now cover placement, residency, scheduling, offload, reserves, and risk acceptance with persistence rules |
| Acceptance was general resource behavior | Acceptance must cover constrained and capable systems, multiple resident models, CPU placement, multi-GPU, scheduling, and no silent substitution |

Why: the previous plan was coherent but could still be read as a conventional single-model loader with better warnings. That reading is now explicitly rejected.

## 3. Owner experience contract

### 3.1 Chat flow

1. User opens Chat.
2. User selects a model from the model dropdown.
3. The selection persists independently of whether the model is resident.
4. User enters a message and presses Send.
5. Howl queues the request.
6. The orchestration engine evaluates current resources and chooses a safe execution strategy.
7. If required, Howl unloads or relocates other workloads, verifies availability as far as supported, and loads the selected model.
8. Howl verifies the model is fully loaded and ready.
9. Only then does Howl dispatch the queued message.

The message must not be sent to an incompletely loaded, failed, or unverified model.

### 3.2 Image flow

1. User switches to Image.
2. The last-used image model remains selected, with settings preserved.
3. User enters a prompt, configures settings, and presses Generate.
4. Howl queues the workflow.
5. The same system-wide engine evaluates resources.
6. If the chat model must leave memory, Howl coordinates the transition safely.
7. Howl verifies release as far as backend and telemetry permit.
8. Howl loads and verifies the image model.
9. Only then does generation begin.

Returning to Chat follows the same process in reverse.

### 3.3 System-wide consistency

The same authority applies to:

- chat inference;
- image generation;
- TTS and speech-related processing;
- music generation;
- 3D mesh generation;
- memory, embedding, retrieval, reranking, and RAG models;
- future model-consuming features.

### 3.4 State distinction that must not be collapsed

| State | Meaning | May survive tab switch? |
|---|---|---|
| Selected | The user's chosen model/settings | Yes |
| Queued | Work is waiting for resources or readiness | Yes, subject to policy |
| Resident | Weights are loaded but not necessarily executing | Yes, if resources permit |
| Executing | The model is actively serving work | Yes, until completion/drain |
| Idle resident | Loaded but not active | Yes |
| Unloading | A transition is in progress | No |
| Unverified | Backend claims unload/load, but evidence is incomplete | Recorded explicitly |
| Failed/interrupted | Transition did not complete | Recorded; no silent fallback |

**Tab switching must not automatically unload.** It may change the visible surface and the default next workload; it does not by itself create a resource transition.

## 4. Unified orchestration responsibilities

The unified authority owns:

1. **Discovery:** CPU, RAM, GPU(s), VRAM, storage, backend, driver, and available health data.
2. **Existing allocations:** usage that existed before Howl, including other applications; the authority must not assume its registry is the whole machine.
3. **Suitability:** weights, runtime overhead, context/KV, workspaces, peak allocation, fragmentation, offload, and backend support.
4. **Placement:** device selection, CPU/GPU choice, layer/offload strategy, multi-GPU assignment, and residency plan.
5. **Admission:** queue, reserve, admit, refuse, or re-plan; no double-spend of the same free bytes.
6. **Transition:** load, verify, drain, reload, relocate, offload, unload, and reconcile.
7. **Execution scheduling:** immediate, concurrent, or sequential execution according to policy and evidence.
8. **Readiness:** dispatch only after the selected model/backend is verified ready for this workload.
9. **Release:** distinguish backend-reported unload from observed physical release and unknown states.
10. **Explanation:** show fit, uncertainty, blockers, transitions, overrides, and failure.
11. **Reconciliation:** detect orphan processes, stale ledger entries, and allocations outside the registry.

The authority must not own prompt content, model quality claims, user state, or RAG corpus authority.

## 5. Residency, concurrency, and scheduling

### 5.1 No arbitrary one-model limit

Howl must not impose a one-model-at-a-time rule. If hardware safely supports a chat model, image model, TTS, music, and 3D model with adequate headroom, the architecture must allow them to remain resident.

The absence of enough resources is a reason to queue or transition, not a reason to ban multi-model residency by design.

### 5.2 Three separate problems

| Problem | Question | Required treatment |
|---|---|---|
| Residency | Which models may remain loaded? | Fit, reserve, peak/workspace, fragmentation, policy |
| Execution concurrency | Which workloads may execute at the same time? | Compute contention, memory peaks, latency, backend support |
| Scheduling | In what order and with what transitions? | Queue policy, priority, fairness, drain, preemption, rollback |

A model may be resident but not executing. Execution may be serialized while multiple models remain resident. Concurrent execution must be justified by measured/supported behavior, not intuition.

### 5.3 Scheduling policy questions

The research program must compare:

- multiple resident models with sequential execution;
- concurrent execution on one GPU;
- concurrent execution across multiple GPUs;
- CPU/GPU mixed execution;
- dynamic load/unload versus keeping models resident;
- priority queues and reservations;
- safe preemption/interruption boundaries;
- hybrid strategies.

The outcome is maximum **useful** capability, not maximum simultaneous activity regardless of stability.

## 6. CPU/RAM versus GPU placement

Placement is a policy decision, not a universal rule.

Howl must consider:

- supported CPU and GPU implementations;
- backend capability and model compatibility;
- RAM versus VRAM capacity and contention;
- latency and throughput;
- device-transfer and synchronization costs;
- hybrid execution and layer offload;
- peak workspace and cache behavior;
- whether GPU acceleration is worth the resource cost for that workload.

CPU execution is not automatically preferable. GPU execution is not automatically preferable. Offloading to RAM is a legitimate placement strategy when supported and appropriate; it must not be an unexplained emergency fallback or silently impose severe degradation.

The user may express a preference; the engine must either honor it or explain why it cannot.

## 7. Multi-GPU as a first-class requirement

Multi-GPU is not a later enhancement or optional curiosity. The design must support, where the backend and runtime actually permit:

- different models on different GPUs;
- one model split across GPUs;
- tensor, pipeline, layer, and other applicable parallelism;
- per-device memory accounting;
- heterogeneous GPU configurations;
- cross-device transfer and synchronization overhead;
- workload placement and balancing;
- vendor-specific capabilities and limits.

The design must not assume identical GPUs, a single vendor, or one universal backend. “Supported” means demonstrated by a named runtime/model/configuration; theoretical possibility is not enough.

## 8. Hardware safety and uncertainty

The authority must account for:

- VRAM and GPU utilization;
- RAM and CPU availability;
- per-device allocations;
- storage/cache requirements;
- weights and runtime overhead;
- KV-cache and context growth;
- peak generation/execution;
- temporary workspaces and auxiliary allocations;
- temperature and other available health data;
- existing allocations outside Howl;
- fragmentation and uncertainty.

The owner's default reserve of at least **10% of currently free VRAM** remains authoritative. It is calculated after existing allocations that existed before Howl loaded models. Research may justify additional peak/fragmentation margins; it may not silently replace the 10% rule.

Telemetry must be lightweight, modular, and vendor-aware. Providers are replaceable. Missing telemetry is explicit, not zero. Conservative estimates and a warned attempt are permitted; missing telemetry must not automatically block all model use.

## 9. Model lifecycle and readiness

### 9.1 Lifecycle states

```text
selected
  → queued
  → admitted / reserved
  → loading
  → verifying
  → ready
  → executing
  → idle-resident
  → draining
  → unloading
  → release-verified | release-partial | release-unknown
  → unloaded

Any transition may become failed / interrupted.
```

### 9.2 Readiness rules

- A load is not ready merely because a process started.
- A backend “loaded” response is not proof that the model can serve the queued workload.
- Verification must match the workload: model identity, backend, context/KV, device placement, and any required auxiliary component.
- A queued request waits for readiness; it is not dispatched to a partial model.
- A failed transition never silently substitutes another model, backend, or RAM placement.
- The selected model remains selected even if temporarily unavailable or not resident.

### 9.3 Release rules

- A successful unload request is not verified physical release.
- Release may be verified, partial, or unknown.
- Runtime caches, driver allocations, in-flight work, shared allocations, KV-cache, and fragmentation may delay or prevent release.
- Unknown release is displayed and retained in the ledger; it is never reported as freed capacity.
- A new admission must reconcile stale/unknown release before claiming the memory.

## 10. Manual control and power-user overrides

Automatic orchestration is the default, not a restriction.

Controls to research and specify:

- preferred GPU/device and model placement;
- CPU/GPU selection where supported;
- residency preference (keep resident versus allow unload);
- automatic transition/unload policy;
- offload preference;
- parallelism and scheduling policy;
- resource reserve;
- risk acceptance and forced attempt.

Overrides must be visible, understandable, persistent where appropriate, and recorded with the exact warning and observed outcome. A forced attempt must not become a silent default.

## 11. Required R-03 research dossier

R-03 is expanded into a substantial program. Each sub-workstream contributes to one dossier with the common evidence rules from the base plan.

| ID | Research focus | Required questions |
|---|---|---|
| **R-03a** | Existing model managers/orchestrators | How do Ollama, LM Studio, llama.cpp-based managers, vLLM, ComfyUI, and comparable systems handle residency, queues, placement, and unload? What fails in practice? |
| **R-03b** | GPU memory allocation/release | How do NVIDIA, AMD, Intel, and relevant runtimes allocate, cache, free, and report memory? What does unload actually release? |
| **R-03c** | Multi-GPU execution/placement | Which parallelism/placement modes are demonstrated by which backends/models? What are transfer, balance, and failure costs? |
| **R-03d** | CPU/GPU placement | When is CPU/RAM execution preferable? What are measured latency, throughput, and memory tradeoffs? |
| **R-03e** | Offload/hybrid execution | Which offload strategies are supported, transparent, stable, and worth their performance cost? |
| **R-03f** | Parallel vs serialized execution | What concurrency is safe and useful on one GPU, multiple GPUs, and mixed CPU/GPU? What are peak/fragmentation effects? |
| **R-03g** | Estimation/admission/headroom | How should weights, KV, workspaces, peaks, and fragmentation be estimated and reserved? |
| **R-03h** | Residency/lifecycle | How should selected, resident, executing, idle, unloading, partial, and unknown states be represented and recovered? |
| **R-03i** | Real-world failures/complaints | What do users report about OOMs, thrashing, stale VRAM, slow reloads, lost selections, silent fallbacks, and multi-GPU failures? |
| **R-03j** | Newer techniques | What newer scheduling, paging, offload, multi-GPU, or memory-management approaches have measured benefits for local creative workloads? |

The dossier must record hardware, model, runtime, configuration, benchmark conditions, and whether evidence is primary documentation, a technical paper, an implementation report, a vendor claim, or a user report.

### 11.1 Research source classes

Start with primary documentation and implementation evidence, not marketing:

- NVIDIA NVML, CUDA memory-management, MPS, and multi-GPU documentation;
- AMD ROCm/HIP/ADL memory and multi-GPU documentation;
- Intel Level Zero/Sysman, oneAPI, and GPU memory documentation;
- llama.cpp model loading, offload, multi-GPU, and backend documentation;
- stable-diffusion.cpp and image-runtime resource behavior;
- TTS/ASR runtime resource behavior;
- Ollama, LM Studio, vLLM, text-generation-webui, ComfyUI, and comparable local orchestrators;
- published systems papers and independent benchmarks with named configurations;
- issue trackers and user reports, clearly labeled anecdotal.

These are research targets, not evidence already established by this amendment.

## 12. Decision-register amendments

The base plan's `DA-03`, `DA-H1`–`DA-H5`, `DA-I4`, and `DA-I8` remain. The following changes and additions apply.

| ID | Decision | Required evidence | Blocks approval | Blocks implementation | Runtime grant |
|---|---|---|---|---|---|
| `DA-H3` | Multi-GPU mode by backend and model; “unsupported” is a per-case finding, not a blanket product choice | R-03c; named runtime/model evidence | Yes | Yes | `OA-09` |
| `DA-H6` | CPU/GPU placement policy and user preference semantics | R-03d; measured tradeoffs | Yes | Yes | `OA-09` |
| `DA-H7` | Residency policy: when models stay resident versus unload | R-03a/h; memory and UX evidence | Yes | Yes | `OA-09` |
| `DA-H8` | Execution scheduling: parallel, serialized, queue, priority, drain, safe preemption | R-03f; stability/peak evidence | Yes | Yes | `OA-09` |
| `DA-H9` | Readiness verification contract per backend/workload | R-03b/h; backend behavior | Yes | Yes | `OA-09` |
| `DA-H10` | Release verification levels and reconciliation rules | R-03b; observed release evidence | Yes | Yes | `OA-09` |
| `DA-H11` | Manual override and persistence contract | R-03d/f/h; UX/security review | Yes | Yes | `OA-09` for behavior |
| `DA-H12` | Additional peak/fragmentation margins beyond the owner's 10% reserve | R-03g; measured evidence | Yes | Yes | `OA-09` |
| `DA-I9` | Unified residency/execution/scheduling contract across all consumers | R-03a–j; integration fixtures | Yes | Yes | `OA-09` |
| `DA-I10` | Model-state persistence and recovery contract | R-03h; restart/crash evidence | Yes | Yes | `OA-09` |
| `DA-I11` | Manual-control surface and override audit contract | R-03d/f/h; UX review | Yes | Yes | `OA-09` |

`DA-H3` changes meaning from “whether multi-GPU is supported at all” to “which placement/parallelism modes are supported for each backend/model, and how does Howl expose them without pretending uniform capability.”

## 13. Integration-map amendment

The base plan's hardware section is expanded as follows.

| Current seam | Current static behavior | Target | Transition/rollback |
|---|---|---|---|
| Resource manager | Advisory registry and hardware snapshot | Unified authority for discovery, placement, residency, admission, scheduling, release | Shadow reports; then one consumer at a time; advisory fallback retained |
| Model selection | Per-chat binding and model settings | Persistent selection independent of residency | Preserve selection when unloaded; no silent substitution |
| Model loading | Loader-driven process start; registration after load; force-load bypass | Queue → admit → load → verify → ready → dispatch | Readiness gate; preserve old path during transition |
| Tab/page switching | No observed automatic unload contract; behavior is not a unified policy | Switching never unloads by itself; policy is explicit | Add no implicit transition; retain manual controls |
| Image loading | Image engine can load independently; fixed estimate registration | Same authority; registered consumer; placement and readiness | Register/observe before enforcement |
| TTS/ASR | Character voice/TTS exists; resource admission incomplete | Same authority; Chat voice is product destination | Add consumer registration before enforcement |
| Memory/RAG embeddings | RAG/embedding models may load through local runtime | Same authority; optional workload class | No separate allocator |
| Multi-GPU | Not a first-class placement contract | Per-device assignments/parallelism where supported | Evidence-gated; no universal claim |
| Release | Fixed 1.5 s wait after unload | Verified/partial/unknown release receipts and reconciliation | Do not remove wait until verified receipts exist |
| Existing allocations | Visible in live free readings but outside registry | Treated as used before reserve | Re-plan; never assume registry completeness |
| Overrides | `force_load` bypasses checks | Visible, persistent, audited override policy | Keep warning; never silent default |

## 14. Failure and recovery requirements

| Failure | Required behavior |
|---|---|
| Selected model unavailable | Keep selection; queue/refuse with explanation; no silent substitution |
| Transition interrupted | Mark failed/interrupted; reconcile; do not report ready |
| Load succeeds, readiness fails | Do not dispatch; retain diagnostics; roll back or unload safely |
| Release unknown/partial | Report unknown/partial; retain ledger; do not spend unverified memory |
| Concurrent work exceeds capacity | Queue/serialize/re-plan; never double-spend |
| Tab switch while executing | Do not kill implicitly; follow drain/preemption policy |
| CPU offload becomes necessary | Disclose placement/performance change; never silently degrade |
| Multi-GPU device disappears | Re-plan/block work; preserve selection; surface error |
| Telemetry unavailable | Conservative estimate, visible uncertainty, allowed warned attempt |
| Existing allocation appears | Treat as used; re-plan; do not evict another application |
| Backend reports loaded/unloaded incorrectly | Trust evidence over claim; reconcile and report |

## 15. Acceptance criteria amendment

All criteria are unvalidated design targets. They require E1–E5 evidence under the base plan's evidence levels.

### 15.1 Constrained systems

- A selected model remains selected after it is unloaded.
- A queued request waits until its selected model is ready.
- Switching tabs does not automatically unload a resident model.
- A new workload can trigger a safe transition when resources are insufficient.
- A transition failure does not silently dispatch to another model or RAM placement.
- Unverified release is not reported as confirmed.
- Existing allocations outside Howl reduce admitted capacity.
- Missing telemetry produces an explicit warned attempt path, not a false zero.

### 15.2 Capable systems

- Multiple models can remain resident when resources permit.
- Concurrent execution is evaluated separately from concurrent residency.
- CPU execution is selectable where it is the better supported strategy.
- Multi-GPU placement works where the backend/model evidence supports it.
- Heterogeneous GPU configurations are handled without assuming one pool.
- Peak/workspace/fragmentation margins are visible and do not silently consume the reserve.
- A model can be idle-resident without being unloaded merely because another tab is active.

### 15.3 Cross-system

- Chat, Image, TTS/ASR, embedding/RAG, Music, 3D, and future consumers all use the same authority.
- No modality tab implements its own conflicting model loader decision.
- Manual overrides are visible, persistent where appropriate, and audited.
- The selected model is never silently replaced.
- Readiness is verified before dispatch.
- Release is verified, partial, or unknown by evidence, never by a fixed sleep.
- Resource contention produces an understandable queue/refusal/transition, not a silent failure.
- Performance claims name hardware, model, runtime, configuration, and workload.

## 16. Dependencies and sequencing

R-03 is now central but is still a **design/research dependency**, not a build order.

```text
R-03a–R-03j
      │
      ▼
DA-03 / DA-H1–DA-H12 / DA-I4 / DA-I9–DA-I11
      │
      ├── memory/RAG workload contracts
      ├── model compatibility (R-04/DA-07)
      ├── Chat/TTS/ASR readiness
      ├── Image readiness
      ├── future Music/3D stage
      └── Android host-capacity policy
      │
      ▼
Unified orchestration contract + acceptance design
      │
      ▼
Owner/architecture approval → implementation authorization
```

Product staging remains Chat → Image → expansion. Residency and scheduling architecture is cross-cutting and must be designed before the first stage's implementation gate, even though some capabilities are only exercised in later stages.

Android may request work, but it never selects a model outside desktop policy and never creates a parallel resource authority.

## 17. Conflicts and unresolved questions

| ID | Conflict or question | Resolution required |
|---|---|---|
| H-01 | Base plan's `DA-H3` allowed “unsupported per backend” to sound like a blanket product choice | Replace with per-backend/model evidence and explicit supported modes |
| H-02 | Base plan distinguished model selection and per-chat binding but not selected/resident/executing | This amendment makes the three states explicit |
| H-03 | Base plan's resource acceptance focused on leases and release, not multi-model residency | Add residency, scheduling, placement, and readiness acceptance |
| H-04 | Base plan treated CPU fallback mainly as degradation | CPU/RAM placement is now a first-class supported strategy |
| H-05 | Base plan's hardware R-03 scope was too narrow | Replace with R-03a–R-03j |
| H-06 | `DA-H2` reserve calibration may conflict with `DA-H12` additional margins | Keep owner's 10% floor; research additional margins separately |
| H-07 | Backend “loaded” and “unloaded” claims may not equal readiness/release | Require evidence-based verification levels |
| H-08 | Multi-GPU support is core, but backend/model support is uneven | Evidence matrix must name supported modes; no universal promise |
| H-09 | Concurrent residency and execution can be confused in UI | Expose separate states and policy |
| H-10 | No R-03 dossier exists yet | This amendment is a research plan, not a finding |

## 18. Authorization state

This amendment:

- does not approve the architecture;
- does not authorize source changes, migrations, model loads, benchmarks, runtime probes, network tests, or pairing;
- does not pass G1–G9;
- does not modify the preserved plan, reconciliation, owner addendum, reviews, or existing design package;
- does not change the product staging or TTS destination.

The next legitimate state is **R-03 research and decision work**, followed by integration/acceptance review and explicit owner authorization.

## 19. Preservation and verification record

| Artifact | Identity |
|---|---|
| Original integration plan | `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf` |
| Original plan size | 53,197 B / 724 lines |
| Original plan sidecar | `docs/architecture/INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md.sha256` |
| This amendment | Recorded in the external amendment sidecar after write |
| Frozen Revision 3 plan | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` |
| Requirements reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` |
| Recovery addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` |
| Independent re-review | `88b12137baca5b27bb733de911c5fb69dc0fb21395123fe0e21dcbcc0fb413af` |

The amendment's own digest is externalized in `AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md.sha256`; no self-hash is embedded.
