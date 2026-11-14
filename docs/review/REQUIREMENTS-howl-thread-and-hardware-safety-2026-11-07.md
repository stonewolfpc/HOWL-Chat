# CROSS-CUTTING ENGINEERING REQUIREMENTS — Thread Safety and Hardware Safety

**Date:** 2026-11-07  
**Status:** RECORDED CROSS-CUTTING ENGINEERING REQUIREMENT — IMPLEMENTATION UNAUTHORIZED  
**Source:** Project-owner instruction recorded in session on 2026-11-07  
**Applies to:** TTS Stage 3 corrective stages **C1–C8** and the eventual implementation  
**Relationship to frozen artifacts:** Additive acceptance overlay. This document does not amend the governing contract, the Stage 3 corrective design, the correction record, the owner-decision record, or any audit.

This is not a performance appendix. Thread safety and hardware/resource safety are peer engineering requirements, not optional tuning notes. They must be carried into every stage's acceptance criteria and into the eventual implementation. Passing a functional test does not establish either requirement.

## 1. Purpose and authority

TH-1, HW-1, and HW-2 are cross-cutting requirements. They apply wherever the implementation shares lifecycle, process, adapter, request, cache, settings, renderer, logging, or status state.

- A stage is not accepted if it does not satisfy the applicable TH-1, HW-1, and HW-2 requirements for that stage.
- A functional success, a performance improvement, or a passing fake-helper test does not override these requirements.
- C8's runtime verification does not retroactively satisfy a requirement that C1–C7 failed to design and test.
- A scheduling deferral, a code comment, or an unrecorded convention is not evidence of compliance.

**Precedence.** The governing contract controls required behavior. These requirements may impose stricter engineering discipline, but they do not authorize a behavior the contract forbids. If a requirement appears to conflict with the contract, implementation must stop and the conflict must be recorded rather than silently resolved in favor of either document.

## 2. Scope

### 2.1 Stages and implementation

The requirements apply to all eight corrective stages:

1. C1 — capability negotiation
2. C2 — status projection
3. C3 — lifecycle and load seam
4. C4 — exit and shutdown consistency
5. C5 — budget, backoff, and idle behavior
6. C6 — identity, recovery, and replacement
7. C7 — production adapter and explicit load
8. C8 — contract closure and runtime verification

They also apply to all code changed by those stages and all code those stages rely on for shared state, process ownership, scheduling, cleanup, or user-visible status.

### 2.2 Application-wide I/O and resource scope

HW-1 is not limited to the TTS subsystem. It applies to new or changed application paths, including:

- settings application and persistence;
- caches and model/assets preparation;
- status and health observation;
- renderer and picker paths;
- logging and diagnostics;
- helper lifecycle and process management;
- temporary/output file handling and cleanup.

The requirement does not authorize an unrelated application-wide refactor. Existing unrelated I/O debt discovered during this work must be recorded separately rather than silently expanded into scope; any path the implementation changes or depends on must meet this overlay.

### 2.3 What this document does not do

This document:

- does not modify the frozen contract or its §17 rules;
- does not modify the Stage 3 corrective design or correction record;
- does not modify the owner-decision record;
- does not resolve D10(a), D10(b), D1a, D9, D1b, D2, D5, or D3;
- does not authorize C1, C2, C3, C4, C5, C6, C7, C8, or implementation;
- does not authorize any build, test, launch, runtime operation, commit, push, or deployment.

## 3. TH-1 — Thread safety

### TH-1.1 Shared-state protection

All shared lifecycle, process, adapter, request, queue, cache, and status state must be protected against concurrent access and race conditions.

- Every shared mutable field must have an explicit synchronization or ownership rule.
- Unsynchronized “fast paths” are not permitted merely because they are expected to be short.
- Snapshot/publication paths must not expose torn, partially updated, or stale-after-terminal state.
- A read path must not hold a lock across blocking I/O, process waits, channel operations, or unbounded work unless the design explicitly proves that the blocking is required and bounded.
- Where multiple locks are used, lock ordering and ownership must be documented and consistently enforced.

### TH-1.2 Defined concurrency behavior

Cancellation, shutdown, load, unload, and synthesis must each have defined concurrency behavior.

The design and tests must state, for every pair or overlap that can occur:

- which operations may run concurrently;
- which operations are serialized, and at what boundary;
- which operation wins when requests race;
- whether a request is rejected, queued, superseded, or completed from an existing state;
- what terminal result is published, and that it is published at most once;
- how cancellation interacts with startup, handshake, loading, active synthesis, and shutdown;
- how a concurrent direct operation differs from the supervisor cancellation path where the contract distinguishes them.

A comment that a race is “unlikely” is not a concurrency specification.

### TH-1.3 Liveness and containment

The implementation must have no deadlocks, unsafe lock ordering, uncontrolled background work, or orphaned workers or processes.

- No goroutine, timer, reader, writer, process, job object, temporary file owner, or cleanup worker may be started without a defined owner and termination path.
- Blocking operations must have bounded waits or an explicit cancellation path.
- Lock ordering must be acyclic and documented where more than one lock can be held.
- A failed operation must not leave a worker, reader, waiter, pipe, process, or temporary file orphaned.
- Shutdown must cause all owned background work to stop or be explicitly accounted for.
- A cancellation that cannot stop a worker must still have a bounded escalation and cleanup path consistent with the contract.

### TH-1.4 Integrity under concurrency

Concurrent requests must not corrupt state or bypass resource limits.

- Admission and process-count checks must be atomic with respect to the resource they protect.
- Queue and concurrency bounds must not be bypassed by a second request, a retry, a status path, or a direct cancellation path.
- A stale completion must not overwrite a newer generation's state or terminal result.
- Duplicate terminal publication, lost terminal publication, send-on-closed-channel, double-close, and cleanup-before-kill ordering must be prevented or explicitly handled.
- Concurrent readers must observe a coherent snapshot, not an interleaving of old and new lifecycle fields.

### TH-1.5 Ownership and shutdown discipline

Shared resources must have one accountable owner and a deterministic release path.

- Process `Wait()` ownership must be single and explicit.
- Kill/wait/cleanup ordering must be deliberate and tested.
- A resource must not be released twice or by two independent owners.
- Shutdown must not race a new operation into starting work after the shutdown decision.
- Background work must be joined or proven terminated before the owner is released.

### TH-1.6 Concurrency evidence

Concurrency claims require more than a successful functional run.

- A race detector or equivalent dynamic race evidence must be used where the language/toolchain supports it.
- Stress and overlap tests must exercise cancellation, shutdown, load/unload, synthesis, retries, and replacement at the boundaries where races are possible.
- Lock-order and ownership claims require static review in addition to tests.
- A test that only exercises one request at a time does not establish thread safety.

## 4. HW-1 — Hardware and resource safety

### HW-1.1 Avoid unnecessary I/O

Avoid unnecessary disk writes and other major I/O across the application.

- Do not add write-on-read behavior, write-per-poll behavior, or a background write cadence without a necessary product or durability reason.
- Do not reread the same model, asset, catalog, or configuration on every request when a bounded cached view is valid.
- Do not flush logs or diagnostics per frame, per status poll, or per operation when a bounded state-change record is sufficient.
- Do not create, rewrite, or delete temporary/output files merely because an operation was observed rather than performed.

### HW-1.2 Necessary writes are permitted

Model-loading writes and other genuinely necessary operations are permitted.

- A write is permitted when it is necessary to materialize or validate a model/asset, publish required output, preserve required settings, or maintain a bounded diagnostic record.
- Necessity must be justified at the call site and bounded in frequency and size.
- A necessary write must not become a repeated write on every request, retry, poll, or inspection.
- Where atomic publication is required by the contract or needed to prevent corruption, it must be used; atomicity is not “unnecessary I/O.”

### HW-1.3 CPU, memory, and background activity

Avoid unnecessary memory allocations, CPU-intensive polling, redundant model discovery, and excessive background activity.

- Do not use a busy loop or tight polling loop in place of a timer, event, condition, or bounded wait.
- Do not perform full model/asset discovery when a cached, invalidated capability snapshot is sufficient.
- Do not allocate unbounded buffers, duplicate large payloads, or retain completed audio/model data without a defined lifetime.
- Do not leave background workers running when their work is complete or their owner is idle.
- Do not add periodic work whose only purpose is to keep a path warm without an explicit product requirement and bounded cost.

### HW-1.4 Bounded concurrency and consumption

Bound concurrency and resource consumption.

Every stage must declare and enforce limits for the resources it introduces or changes, including as applicable:

- worker/goroutine counts;
- queue depth and request concurrency;
- retry and restart counts;
- timer and background-task counts;
- process count;
- memory and buffer sizes;
- open files, pipes, handles, and job objects;
- temporary/output storage;
- log volume and retention.

“Unbounded” is not an acceptable default. A limit may be provisional, but it must be explicit, testable, and consistent with the frozen contract.

### HW-1.5 Hardware-aware scheduling

Use hardware-aware scheduling and workload management so TTS does not monopolize CPU, memory, GPU, disk, or other resources.

- Concurrency must be appropriate to available CPU capacity and must leave headroom for the application and the user.
- A TTS workload must not assume a high-end machine, a dedicated GPU, or unlimited memory.
- Where platform scheduling or priority controls are used, they must be bounded, reversible, and must not degrade foreground responsiveness.
- Resource-intensive preparation or inference must not be started from inspection, status, startup, or catalog paths unless the contract explicitly permits that start.
- The implementation must avoid creating a thundering herd of helpers, workers, timers, or retries when several requests arrive together.

### HW-1.6 Cross-hardware safety

Prefer techniques that work safely across a range of hardware configurations, not just high-end machines.

- The design must state its assumptions about CPU cores, memory, disk, and GPU availability.
- Low-core and low-memory behavior must be considered, not only the developer workstation.
- A missing or weaker GPU must not turn into an unsafe retry loop, an unbounded fallback, or a silent correctness change.
- Degraded performance is acceptable only when behavior, cleanup, and lifecycle correctness remain intact.

### HW-1.7 Observable contention and degradation

Resource contention and degradation must be observable and handled deliberately.

- The implementation must expose bounded, low-cost signals for queueing, rejection, retry/backoff, process count, and resource pressure where those conditions affect behavior.
- Observability must not itself violate HW-1.1 or HW-1.3; bounded counters and state changes are preferred over high-frequency logging.
- Degradation must have a defined response: wait within a bound, queue within a bound, reject with the contract-permitted result, or reduce nonessential work.
- Silent spinning, silent dropping, and unbounded retry are not deliberate handling.

## 5. HW-2 — Performance must not compromise safety

### HW-2.1 Safety takes precedence

Never sacrifice synchronization, data integrity, process cleanup, or lifecycle correctness just to improve speed.

Performance work must not introduce:

- an unsynchronized fast path;
- a fire-and-forget worker without an owner or termination path;
- a detached process without kill/wait/cleanup ownership;
- an unbounded queue, retry, timer, or allocation;
- a skipped cleanup or best-effort release justified only by latency;
- a status or cache shortcut that can expose stale terminal state;
- a resource-limit bypass justified by expected load.

### HW-2.2 Deliberate degradation

Resource contention and degradation must be handled deliberately, not hidden.

- A slow path must remain a correct path.
- Backpressure, rejection, queueing, or reduced nonessential work must be explicit and bounded.
- A performance change must preserve the same terminal results, cleanup ordering, and ownership guarantees as the safe path.
- If the safe behavior cannot meet a performance target, the target is revised or the work is escalated; safety is not silently weakened.

### HW-2.3 Performance evidence

A performance claim must include safety evidence.

- Before/after measurements must state the workload, hardware profile, and resource dimensions.
- A speedup measured only on a high-end machine does not establish cross-hardware safety.
- A performance optimization must include the concurrency, cleanup, and resource-limit tests affected by the change.
- A regression in safety evidence blocks acceptance regardless of the speed result.

## 6. Stage acceptance rules

Every stage acceptance must include the following. These rules are cumulative; a later stage does not excuse an earlier omission.

| Rule | Requirement |
|---|---|
| **AC-1 — Concurrency model** | The stage documents the shared state it touches, the operations that can overlap, the serialization boundary, and the winner/terminal rule for each relevant race. |
| **AC-2 — Lock and ownership review** | The stage records lock ordering, single-owner resources, blocking-under-lock exclusions, and background-work ownership. Static review is required where tests cannot prove the property. |
| **AC-3 — Dynamic concurrency evidence** | The stage runs race/stress/overlap tests for its affected operations, including cancellation, shutdown, load/unload, retry, and replacement where applicable. Evidence must be repeatable, not a single happy-path run. |
| **AC-4 — Resource budget** | The stage declares bounded limits for its workers, queues, retries, timers, processes, buffers, handles, I/O, and storage, and records the measurements or checks that support them. |
| **AC-5 — Safety under load** | The stage demonstrates that concurrent requests and degraded resources do not corrupt state, bypass limits, leak workers/processes, or lose/duplicate terminal results. |
| **AC-6 — No hidden background work** | Every started goroutine, timer, reader, process, and cleanup task has an owner, a bound, and a termination path. |
| **AC-7 — No stage substitution** | C8's final runtime verification may aggregate evidence, but it cannot substitute for the stage-level evidence C1–C7 must produce. |

A stage that has functional tests but no applicable TH/HW evidence is **not accepted**.

## 7. C1–C8 acceptance overlay

The following table is binding for stage acceptance. It supplements, and does not replace, each stage's functional acceptance criteria in the frozen design.

| Stage | TH-1 acceptance focus | HW-1 / HW-2 acceptance focus | Required evidence |
|---|---|---|---|
| **C1 — Capability negotiation** | Capability snapshot publication is atomic and single-writer; concurrent negotiation, gating, and reads do not observe torn or stale fields; handshake completion and failure have one terminal outcome. | No spawn, model load, or disk write from negotiation or inspection; parsing and descriptor storage are bounded; no redundant discovery per operation. | Race/stress tests for concurrent handshake and capability reads; invariant check that unknown/absent paths do not start work; allocation/limit review. |
| **C2 — Status projection** | Multi-field status is a coherent snapshot; lifecycle transitions cannot produce mixed old/new fields; status reads do not hold locks across I/O or block transitions indefinitely. | Status is bounded-cost and performs no model discovery, file read, or high-frequency log write; buffers and counters are bounded. | Concurrent status-vs-transition tests; no-torn-snapshot assertions; review of blocking calls under locks; I/O trace or design proof. |
| **C3 — Lifecycle and load seam** | Load/unload/cancel/synthesis have defined overlap and winner rules; load is single-flight where required; lock order between adapter, client, and lifecycle state is documented; no orphaned worker on failure. | No eager spawn or model load from inspection/status/startup; no duplicate load; waits are bounded; no write cadence added by lifecycle bookkeeping. | Overlap tests for load/unload/cancel; failure-path ownership checks; lock-order review; duplicate-load and orphan checks. |
| **C4 — Exit and shutdown consistency** | Cancel/shutdown/kill/wait races have one terminal result; single `Wait()` ownership holds; no double-close, send-on-closed, or cleanup-before-kill; shutdown prevents new work from starting. | Grace and shutdown waits are bounded; no busy wait; no unnecessary file rewrite during cleanup; process and handle release is prompt. | Race tests for shutdown vs active work; single-wait and double-terminal assertions; goroutine/process/handle leak checks; cleanup-ordering review. |
| **C5 — Budget, backoff, and idle** | Restart ledger, backoff timers, idle deadline, and new requests have defined races; no lost wakeup, double close, timer leak, or deadlock between timer and request paths. | Timers/events replace CPU polling; restart count and backoff are bounded; no disk churn per tick; idle close actually releases resources. | Timer-vs-request overlap tests; restart/backoff bound tests; idle-close resource check; race/stress evidence. |
| **C6 — Identity, recovery, and replacement** | Identity transitions are serialized; replacement vs pending requests has a defined winner; stale completions cannot overwrite the new generation; cap admission cannot be bypassed. | Old resources are released before/while replacement without overlapping unbounded model memory; no redundant model-load writes; replacement is bounded. | Concurrent identity/replacement tests; stale-generation assertions; resource release/overlap review; process-count checks. |
| **C7 — Production adapter and explicit load** | Adapter methods are safe under concurrent calls; direct cancel/unload vs synthesis is defined; cap admission is atomic; no path bypasses the limit; explicit load cannot race shutdown. | Adapter is default-off; no startup load; catalog/discovery work is cached and bounded; no unnecessary renderer/status/disk I/O; queues and helpers are bounded. | Adapter concurrency tests; cap-admission race tests; default-off and no-eager-load checks; resource budget and I/O review; caller-map test. |
| **C8 — Runtime verification** | Race/stress/soak coverage for cancellation, shutdown, load/unload, retry, replacement, and cap admission; lock-order, goroutine, process, handle, and temporary-file leak checks; clean repeated runs. | Measure CPU, memory, disk I/O, GPU/accelerator use where applicable, handles, and process count at idle and under load; include a low-core/low-memory profile; show contention observability and no safety regression. | Repeatable runtime report with hardware profile, workload, measurements, and failure-path results; cannot be satisfied by fake-helper-only evidence where the claim is about real hardware/resource behavior. |

## 8. Evidence standard

The following evidence rules apply to every stage:

1. **Static evidence** — lock inventory, ownership map, blocking-call review, and lock-order rationale.
2. **Dynamic evidence** — race/stress/overlap tests that exercise the actual concurrency boundary, not only single-request behavior.
3. **Resource evidence** — declared bounds plus measurements or checks for the resources that stage affects.
4. **Failure evidence** — cancellation, shutdown, timeout, start failure, replacement, and cap rejection paths, not only success paths.
5. **Repeatability** — a single passing run is not sufficient for a race or leak claim; the evidence must be repeatable and its workload stated.
6. **Scope honesty** — fake-helper evidence cannot establish real hardware behavior, and a code review cannot establish a dynamic absence of races by itself.

A stage may use provisional numeric limits where the contract already labels a value provisional, but the limit must still be explicit bounded and testable.

## 9. Application-wide I/O policy

The application-wide I/O requirement is deliberately broader than the TTS helper:

- **Settings:** do not rewrite the settings file on read, status refresh, or an unchanged apply.
- **Caches:** do not rescan or rewrite a valid cache on every request; invalidate deliberately.
- **Logging:** do not emit per-frame, per-poll, or per-operation disk logs when a bounded state-change record is sufficient.
- **Status/health:** do not perform model discovery or file reads merely to answer a status request.
- **Renderer/catalog:** do not start a helper or load a model from an inspection path unless the contract and owner decision permit it.
- **Output/temporary files:** write only what is required, publish atomically where needed, and clean up deterministically.
- **Model loading:** materialization and validation writes are permitted when necessary; they must not become a per-request or per-retry rewrite.

An unrelated legacy I/O problem discovered during this work is recorded as a separate finding. It is not used to justify expanding the current stage, and it is not silently fixed outside scope.

## 10. Observability and degradation

Resource contention and degradation must be observable without creating a new resource problem.

- Prefer bounded counters and state changes over high-frequency logs.
- Record queue depth, rejection, retry/backoff, process count, and resource-pressure indicators where those conditions affect behavior.
- A degradation response must be one of: wait within a bound, queue within a bound, reject with a contract-permitted result, or reduce nonessential work.
- Observability must not hold lifecycle locks across I/O or allocate unbounded per-request history.
- A health/status response must not become a hidden model-discovery or disk-scan operation.

## 11. Conflict, exception, and waiver handling

- **Contract conflict:** if a TH/HW requirement cannot be met without violating the frozen contract, stop and record the conflict. Do not silently weaken either document.
- **Stricter requirement:** a stricter TH/HW rule may be applied when it is compatible with the contract.
- **No performance waiver:** TH-1, HW-1, and HW-2 may not be waived merely to improve latency, throughput, or benchmark results.
- **No silent exception:** any approved exception must state its scope, duration, owner, and the safety evidence that remains in force.
- **No implied authorization:** an exception, a stage acceptance, or this requirements record does not authorize implementation, C8, or deployment.

## 12. Relationship to pending owner decisions

This document records engineering requirements only. It does not resolve any pending owner question.

- D10(a) and D10(b) remain pending in the owner-decision instrument.
- D1a, D9, D1b, D2, D5, and D3 remain pending.
- A conversational preference for a D10(a) result is not a recorded owner decision unless it is completed in the owner-decision instrument; this document neither adopts nor resolves it.
- C1 and C2 remain unauthorized; all later stages remain unauthorized; C8 still requires separate runtime-execution authorization.

## 13. Completion checklist

Before a stage is proposed for acceptance, the reviewer must be able to answer yes to all applicable items:

- [ ] The stage's shared state and ownership map is documented.
- [ ] Every overlap of cancellation, shutdown, load, unload, and synthesis has a defined winner and terminal result.
- [ ] Lock ordering is documented and no blocking I/O occurs under a lifecycle lock without a bounded justification.
- [ ] Every background worker, timer, process, and cleanup task has an owner and termination path.
- [ ] Concurrent requests cannot corrupt state or bypass process/queue/resource limits.
- [ ] The stage's I/O is necessary, bounded, and free of write-on-read or write-per-poll behavior.
- [ ] Worker, queue, retry, timer, process, memory, handle, and storage bounds are explicit.
- [ ] The workload is safe on a low-core/low-memory profile, not only a high-end machine.
- [ ] Contention and degradation are observable and have a deliberate response.
- [ ] No speed optimization weakens synchronization, integrity, cleanup, or lifecycle correctness.
- [ ] The stage has race/stress/overlap and resource evidence, not only functional tests.
- [ ] C8 is not being used to substitute for missing C1–C7 evidence.

A checked box is a claim requiring evidence, not a waiver.
