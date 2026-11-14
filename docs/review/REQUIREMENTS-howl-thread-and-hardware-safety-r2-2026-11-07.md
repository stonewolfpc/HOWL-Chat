# CROSS-CUTTING ENGINEERING REQUIREMENTS — Thread Safety, Hardware Safety, and Performance Preservation

**Revision:** r2 corrective  
**Date:** 2026-11-07  
**Status:** PROPOSED CORRECTIVE REVISION — NOT ADOPTED — NOT BINDING  
**Audited predecessor:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md` (SHA-256 `ee38f0310042188a6821f8632e9e07bfa60dcac454a71ff2d972bb47439ac0e5`), preserved unchanged  
**Pre-adoption gate:** Independent Mimo review pending; owner adoption pending  
**Applies to:** TTS Stage 3 corrective stages **C1–C8** and the eventual implementation  
**Authority:** Additive engineering requirement proposal. It does not amend the governing contract, the Stage 3 corrective design, the correction record, the owner-decision record, or any audit. It becomes binding only through the separate owner-adoption record in §2.

This revision is a response to the read-only audit findings F1–F10 and the project-owner requirement to preserve normal-operation performance. It is not an adoption, an authorization, a stage acceptance, or a decision on any pending owner question.

## 1. Corrective basis

| Audit finding | Resolution in this revision |
|---|---|
| **F1 — authority** | §2 defines a separate, hash-bound owner-adoption mechanism. This document is explicitly non-binding until adoption. |
| **F2 — C1 handshake contradiction** | C1’s no-spawn rule is scoped to discovery of optional capabilities; a contract-authorized handshake start remains permitted. See §8, C1. |
| **F3 — C6 overlap** | Replacement requires shutdown, termination if necessary, and observed `Wait()` before any new process starts; no concurrent old/new residency. See §8, C6. |
| **F4 — timeout/helper-exit races** | Timeout and unexpected helper exit are first-class TH-1.2 events and C4/C8 evidence. See §4.2, §8 C4, and §8 C8. |
| **F5 — 120 s idle grace** | The idle timer and helper may remain through the contract-defined 120 s grace; premature shutdown is not required. See §5.3 and §8 C5. |
| **F6 — app-wide I/O evidence** | §9 defines a mandatory I/O impact record for every stage and every touched application path. See AC-7 and §8. |
| **F7 — measurable budgets** | §10 defines a resource-budget record; “bounded” requires a numeric value or explicit contract-provisional value, a measurement method, and an acceptance threshold. |
| **F8 — exception loophole** | §12 prohibits performance waivers of TH-1, HW-2, contract cleanup/lifecycle/data-integrity requirements, and process ordering. Any HW-1 budget exception requires a recorded owner decision with scope, expiry, and safety evidence. |
| **F9 — observability surface** | §11 makes bounded internal diagnostics the default. No renderer, public API, or status-shape change is authorized without separate contract/owner authority. |
| **F10 — C8 pass/fail** | §6.4 and §8 C8 require a pre-execution test plan with hardware profile, workload, duration/iterations, budgets, and thresholds, including the performance-regression check. |
| **New requirement — performance preservation** | HW-1.8 and HW-2.4 add normal-operation performance preservation, contention-only adaptive controls, separation of instrumentation overhead from contention degradation, and a measurable regression check. |

No finding justifies abandoning the safety requirements. The predecessor remains the audited revision; this file is the proposed corrective revision for independent review.

## 2. Authority and owner adoption

### 2.1 Current status

This document is **PROPOSED — NOT BINDING**. Until the owner completes a separate adoption record, it may be used as a review checklist but cannot be used to accept or reject a stage, authorize implementation, or waive any contract requirement.

The governing contract remains controlling. This document does not amend §17 or any adopted clarification. If a requirement conflicts with the contract, the conflict must be recorded and escalated; neither document may be silently overridden.

### 2.2 Separate adoption record required

A document cannot adopt itself, and inserting this file’s final hash into this file would change the hash. Adoption must therefore be recorded separately and must bind the exact artifact. The adoption record must contain:

- owner identity or accountable role;
- adoption date;
- exact path and SHA-256 of this file;
- adopted scope, including C1–C8 and the eventual implementation;
- confirmation that the governing contract and frozen design inputs were not amended;
- the independent Mimo review record and its verified hash;
- any approved exceptions, with scope, duration, and safety evidence;
- an explicit statement that adoption does not authorize any stage or implementation.

### 2.3 Adoption block — left blank

**Owner:**  
**Adoption date:**  
**Artifact path:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r2-2026-11-07.md`  
**Artifact SHA-256:**  
**Adopted scope:**  
**Mimo review record / hash:**  
**Contract-amendment confirmation:**  
**Approved exceptions:**  
**Decision:** [ ] Adopt  [ ] Reject  [ ] Return for correction

### 2.4 Pre-adoption rule

Owner adoption must not occur before the independent Mimo review is complete and recorded. A Mimo finding that identifies a safety, authority, or contract contradiction reopens this revision rather than being treated as an adoption prerequisite waived by schedule.

No adoption, review, or stage acceptance authorizes C1, C2, C3, C4, C5, C6, C7, C8, implementation, or runtime execution.

## 3. Scope

### 3.1 Stages and implementation

This revision applies to all eight corrective stages:

1. C1 — capability negotiation
2. C2 — status projection
3. C3 — lifecycle and load seam
4. C4 — exit and shutdown consistency
5. C5 — budget, backoff, and idle behavior
6. C6 — identity, recovery, and replacement
7. C7 — production adapter and explicit load
8. C8 — contract closure and runtime verification

It also applies to code changed by those stages and code directly integrated with their shared state, process ownership, scheduling, cleanup, or user-visible status.

“All code those stages rely on” is limited to code the stage changes or directly integrates. It does not authorize an audit or retrofit of unrelated legacy code, third-party internals, or the entire application. Unrelated I/O debt discovered during the work is recorded separately under §9.4.

### 3.2 Application-wide I/O scope

HW-1 applies to every new or changed application path, not only the TTS helper. This includes settings, caches, model/asset preparation, status and health observation, renderer and picker paths, logging and diagnostics, helper lifecycle, and temporary/output files.

“Application-wide” means all paths touched or directly depended on by the implementation. It does not mean that every pre-existing unrelated write in the application must be fixed by this program. The distinction must be recorded per stage in the I/O impact record (§9).

### 3.3 What this document does not do

This revision does not:

- amend the frozen contract or its §17 rules;
- amend the Stage 3 corrective design or correction record;
- amend the owner-decision record;
- resolve D10(a), D10(b), D1a, D9, D1b, D2, D5, or D3;
- accept or authorize C1–C8;
- authorize implementation, build, test, launch, benchmark, runtime operation, commit, push, or deployment.

## 4. TH-1 — Thread safety

### 4.1 TH-1.1 Shared-state protection

All shared lifecycle, process, adapter, request, queue, cache, and status state must be protected against concurrent access and race conditions.

- Every shared mutable field must have an explicit synchronization or single-owner rule.
- Unsynchronized fast paths are not permitted merely because they are expected to be short.
- Snapshot and publication paths must not expose torn, partially updated, or stale-after-terminal state.
- A read path must not hold a lock across blocking I/O, process waits, channel operations, or unbounded work unless the design explicitly proves the blocking is required and bounded.
- Lock ordering and ownership must be documented where more than one lock can be held.

### 4.2 TH-1.2 Defined concurrency behavior

Cancellation, shutdown, timeout, load, unload, synthesis, and unexpected helper exit must each have defined concurrency behavior.

The design and tests must state, for every overlap that can occur:

- which operations may run concurrently;
- which operations are serialized, and at what boundary;
- which operation wins when requests race;
- whether a request is rejected, queued, superseded, or completed from an existing state;
- what terminal result is published, and that it is published at most once;
- how cancellation and timeout interact with startup, handshake, loading, active synthesis, and shutdown;
- how an unexpected `Wait()`/EOF exit interacts with an active request, cancellation, timeout, and cleanup;
- how a late terminal response is discarded or rejected after timeout, cancellation, or retirement;
- how a concurrent direct operation differs from the supervisor cancellation path where the contract distinguishes them.

A comment that a race is “unlikely” is not a concurrency specification.

### 4.3 TH-1.3 Liveness and containment

The implementation must have no known deadlock, unsafe lock ordering, uncontrolled background work, or orphaned workers or processes.

- No goroutine, timer, reader, writer, process, job object, temporary-file owner, or cleanup worker may start without a defined owner and termination path.
- Blocking operations must have bounded waits or an explicit cancellation path.
- Lock ordering must be acyclic and documented where more than one lock can be held.
- A failed operation must not leave a worker, reader, waiter, pipe, process, or temporary file orphaned.
- Shutdown must cause all owned background work to stop or be explicitly accounted for.
- A cancellation that cannot stop a worker must have a bounded escalation and cleanup path consistent with the contract.
- The contract-defined idle grace is not an orphan: a timer or helper waiting for the 120 s idle-close trigger is permitted until that trigger or an earlier authorized shutdown.

### 4.4 TH-1.4 Integrity under concurrency

Concurrent requests must not corrupt state or bypass resource limits.

- Admission and process-count checks must be atomic with respect to the resource they protect.
- Queue and concurrency bounds must not be bypassed by a second request, a retry, a status path, or a direct cancellation path.
- A stale completion must not overwrite a newer generation's state or terminal result.
- Duplicate terminal publication, lost terminal publication, send-on-closed-channel, double-close, and cleanup-before-kill ordering must be prevented or explicitly handled.
- Concurrent readers must observe a coherent snapshot, not an interleaving of old and new lifecycle fields.
- Timeout and unexpected exit must not race a success into acceptance after cancellation or retirement has won.

### 4.5 TH-1.5 Ownership and shutdown discipline

Shared resources must have one accountable owner and a deterministic release path.

- Process `Wait()` ownership must be single and explicit.
- Kill/wait/cleanup ordering must be deliberate and tested.
- A resource must not be released twice or by two independent owners.
- Shutdown must not race a new operation into starting work after the shutdown decision.
- Background work must be joined or proven terminated before the owner is released.

### 4.6 TH-1.6 Concurrency evidence

Concurrency claims require more than a successful functional run.

- A race detector or equivalent dynamic race evidence must be used where the language/toolchain supports it.
- Stress and overlap tests must exercise cancellation, shutdown, timeout, load/unload, synthesis, retries, replacement, and unexpected helper exit at the boundaries where races are possible.
- Tests must include timeout versus cancellation, timeout versus exit, cancellation versus exit, and late terminal versus retired request.
- Lock-order and ownership claims require static review in addition to tests.
- A test that exercises one request at a time does not establish thread safety.

### 4.7 TH-1.7 Applicability

Every stage must record a TH-1 applicability decision. “Not applicable” is permitted only when the stage has no shared state or concurrency surface, and the decision must be explicitly reviewed. A blanket “not applicable” statement is not evidence.

## 5. HW-1 — Hardware and resource safety

### 5.1 HW-1.1 Avoid unnecessary I/O

Avoid unnecessary disk writes and other major I/O across the application.

- Do not add write-on-read, write-per-poll, or background write cadence without a necessary product or durability reason.
- Do not reread the same model, asset, catalog, or configuration on every request when a bounded, valid cached view is sufficient.
- Do not flush logs or diagnostics per frame, per status poll, or per operation when a bounded state-change record is sufficient.
- Do not create, rewrite, or delete temporary/output files merely because an operation was observed rather than performed.
- “Major I/O” and “unnecessary” must be interpreted through the I/O impact record in §9, not by an undocumented reviewer intuition.

### 5.2 HW-1.2 Necessary writes are permitted

Model-loading writes and other genuinely necessary operations are permitted.

- A write is permitted when necessary to materialize or validate a model/asset, publish required output, preserve required settings, or maintain a bounded diagnostic record.
- Necessity must be justified at the call site and bounded in frequency and size.
- A necessary write must not become a repeated write on every request, retry, poll, or inspection.
- Where atomic publication is required by the contract or needed to prevent corruption, it must be used; atomicity is not unnecessary I/O.

### 5.3 HW-1.3 CPU, memory, and background activity

Avoid unnecessary memory allocations, CPU-intensive polling, redundant model discovery, and excessive background activity.

- Do not use a busy loop or tight polling loop in place of a timer, event, condition, or bounded wait.
- Do not perform full model/asset discovery when a cached, invalidated capability snapshot is sufficient.
- Do not allocate unbounded buffers, duplicate large payloads, or retain completed audio/model data without a defined lifetime.
- Do not leave completed background work running after its owner is finished.
- Do not add periodic work whose only purpose is to keep a path warm without an explicit product requirement and bounded cost.
- The contract-defined 120 s idle grace is an explicit exception: a helper or idle timer may remain until the idle-close trigger, and the timer must not be treated as an orphan or a violation of this rule.

### 5.4 HW-1.4 Bounded concurrency and consumption

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

“Unbounded” is not an acceptable default. A limit must be explicit, testable, and recorded in the resource-budget record in §10. A limit may use a contract-provisional value, but the source of that value must be named.

### 5.5 HW-1.5 Hardware-aware scheduling

Use hardware-aware scheduling and workload management so TTS does not monopolize CPU, memory, GPU, disk, or other resources.

- Concurrency must be appropriate to available CPU capacity and must leave headroom for the application and the user.
- A TTS workload must not assume a high-end machine, a dedicated GPU, or unlimited memory.
- Where platform scheduling or priority controls are used, they must be bounded, reversible, and must not degrade foreground responsiveness.
- Resource-intensive preparation or inference must not start from inspection, status, startup, or catalog paths unless the contract explicitly permits that start.
- The implementation must avoid a thundering herd of helpers, workers, timers, or retries when several requests arrive together.

### 5.6 HW-1.6 Cross-hardware safety

Prefer techniques that work safely across a range of hardware configurations, not just high-end machines.

- The design must state its assumptions about CPU cores, memory, disk, and GPU availability.
- Low-core and low-memory behavior must be considered, not only the developer workstation.
- A missing or weaker GPU must not turn into an unsafe retry loop, an unbounded fallback, or a silent correctness change.
- Degraded performance is acceptable only when behavior, cleanup, and lifecycle correctness remain intact.

### 5.7 HW-1.7 Observable contention and degradation

Resource contention and degradation must be observable and handled deliberately.

- The implementation must expose bounded, low-cost internal signals for queueing, rejection, retry/backoff, process count, and resource pressure where those conditions affect behavior.
- Observability must not violate HW-1.1 or HW-1.3; bounded counters and state changes are preferred over high-frequency logging.
- Degradation must have a defined response: wait within a bound, queue within a bound, reject with a contract-permitted result, or reduce nonessential work.
- Silent spinning, silent dropping, and unbounded retry are not deliberate handling.
- The observation surface follows §11: internal diagnostics by default; no public status, renderer, or API field is added without separate authority.

### 5.8 HW-1.8 Performance preservation during normal operation

Hardware protection must preserve available performance during normal operation.

- No arbitrary CPU or GPU utilization ceiling may be imposed during normal operation.
- No unnecessary throttling, sleep, delay, or concurrency reduction may be applied when no contention is measured.
- A contract-fixed process cap or a justified bounded concurrency limit is not an “arbitrary ceiling”; it must be recorded, justified, and must not be lowered further without contention evidence.
- Adaptive controls must activate in response to actual, measured contention, must be reversible, and must deactivate when contention clears.
- A static blanket throttle applied to normal operation is not an acceptable adaptive control.
- The overhead of safety instrumentation and verification must be measured separately from contention-induced degradation.
- Safety must never be sacrificed to improve benchmark results, and a benchmark result must never be used to justify weakening synchronization, cleanup, or lifecycle correctness.

## 6. HW-2 — Performance must not compromise safety

### 6.1 HW-2.1 Safety takes precedence

Never sacrifice synchronization, data integrity, process cleanup, or lifecycle correctness just to improve speed.

Performance work must not introduce:

- an unsynchronized fast path;
- a fire-and-forget worker without an owner or termination path;
- a detached process without kill/wait/cleanup ownership;
- an unbounded queue, retry, timer, or allocation;
- a skipped cleanup or best-effort release justified only by latency;
- a status or cache shortcut that can expose stale terminal state;
- a resource-limit bypass justified by expected load;
- an unmeasured permanent throttle justified only by “safety.”

### 6.2 HW-2.2 Deliberate degradation

Resource contention and degradation must be handled deliberately, not hidden.

- A slow path must remain a correct path.
- Backpressure, rejection, queueing, or reduced nonessential work must be explicit and bounded.
- A performance change must preserve the same terminal results, cleanup ordering, and ownership guarantees as the safe path.
- If the safe behavior cannot meet a performance target, the target is revised or the work is escalated; safety is not silently weakened.

### 6.3 HW-2.3 Performance evidence

A performance claim must include safety evidence.

- Before/after measurements must state the workload, hardware profile, and resource dimensions.
- A speedup measured only on a high-end machine does not establish cross-hardware safety.
- A performance optimization must include the concurrency, cleanup, and resource-limit tests affected by the change.
- A regression in safety evidence blocks acceptance regardless of the speed result.

### 6.4 HW-2.4 Measurable performance-regression check

A performance-regression check is mandatory for C8 and for any stage that changes scheduling, concurrency, throttling, caching, or normal-operation work. The check must use the same declared hardware profile and workload for baseline and candidate, with warm-up and repeated runs.

| Check | Measurement | Pass condition |
|---|---|---|
| **PR-1 — Normal-operation overhead** | No injected contention; baseline versus candidate; latency (median and p95), throughput, CPU, memory, disk I/O, handles, and process count where applicable | Candidate does not exceed the approved regression threshold; no arbitrary CPU/GPU ceiling or unnecessary throttle is active |
| **PR-2 — Contention degradation** | Same workload under a declared contention profile; baseline versus candidate | Degradation is bounded by the approved threshold, and safety behavior remains correct |
| **PR-3 — Adaptive control** | Activation and deactivation under measured contention | Controls activate only after actual contention, are reversible, and deactivate when contention clears; no permanent normal-operation throttle |
| **PR-4 — Safety regression** | Race, timeout, cancellation, exit, cleanup, and resource-limit checks under the same workload | Any safety regression blocks acceptance regardless of PR-1/PR-2 results |

The numeric regression thresholds, contention profile, run count, and warm-up period must be recorded and approved before C8 execution. This document does not invent them. A report that presents measurements without a pass/fail comparison is not C8 evidence.

## 7. Stage acceptance rules

Every stage acceptance must include the following. These rules are cumulative; a later stage does not excuse an earlier omission.

| Rule | Requirement |
|---|---|
| **AC-1 — Concurrency model** | The stage documents shared state, overlapping operations, serialization boundaries, and winner/terminal rules for each relevant race. |
| **AC-2 — Lock and ownership review** | The stage records lock ordering, single-owner resources, blocking-under-lock exclusions, and background-work ownership. |
| **AC-3 — Dynamic concurrency evidence** | The stage runs race/stress/overlap tests for cancellation, shutdown, timeout, load/unload, retry, replacement, and unexpected exit where applicable. |
| **AC-4 — Resource budget** | The stage records every applicable limit in the §10 budget record, with a numeric value or contract-provisional value, measurement method, and acceptance threshold. |
| **AC-5 — Safety under load** | The stage demonstrates that concurrent requests and degraded resources do not corrupt state, bypass limits, leak workers/processes, or lose/duplicate terminal results. |
| **AC-6 — No hidden background work** | Every started goroutine, timer, reader, process, and cleanup task has an owner, a bound, and a termination path, except the contract-defined idle grace. |
| **AC-7 — Application I/O evidence** | Every stage records its I/O impact under §9, including explicit N/A for untouched paths. |
| **AC-8 — Performance preservation** | The stage records whether its changes affect normal-operation performance; C8 applies PR-1 through PR-4. |
| **AC-9 — Applicability review** | A TH-1 or HW-1 “not applicable” determination requires an explicit reviewed entry; blanket N/A is not accepted. |
| **AC-10 — No stage substitution** | C8 may aggregate evidence but cannot substitute for missing C1–C7 evidence. |

A stage that has functional tests but no applicable TH/HW evidence is not accepted. A stage with no declared resource budget is not accepted.

## 8. C1–C8 acceptance overlay

The following table supplements—and does not replace—each stage's frozen functional acceptance criteria and every pending owner-decision gate. The D10(a)/D10(b), D1a, D9, D1b, D2, D5, and D3 gates remain in force.

| Stage | TH-1 acceptance focus | HW-1 / HW-2 / performance focus | Required evidence |
|---|---|---|---|
| **C1 — Capability negotiation** | Capability snapshot publication is atomic and single-writer; concurrent negotiation, gating, and reads do not observe torn or stale fields; handshake completion and failure have one terminal outcome. | No spawn solely to discover an optional capability; a contract-authorized handshake start remains permitted. No model load or disk write from capability negotiation or inspection. Parsing and descriptor storage are bounded. | Race/stress tests for concurrent handshake and capability reads; structural no-frame/no-spawn checks for optional-operation discovery; authorized-handshake start evidence; allocation/limit review; I/O impact record (likely N/A). |
| **C2 — Status projection** | Multi-field status is a coherent snapshot; lifecycle transitions cannot produce mixed old/new fields; status reads do not hold locks across I/O or block transitions indefinitely. | Status is bounded-cost and performs no model discovery, file read, or high-frequency log write; buffers and counters are bounded; observability is internal §11 diagnostics only. | Concurrent status-vs-transition tests; no-torn-snapshot assertions; blocking-call review; I/O impact record; evidence that no public status/renderer/API field was added. |
| **C3 — Lifecycle and load seam** | Load/unload/cancel/synthesis have defined overlap and winner rules; one active load; lock order between adapter, client, and lifecycle state is documented; no orphaned worker on failure. | No eager spawn or model load from inspection/status/startup; no duplicate load; waits are bounded; no write cadence added by lifecycle bookkeeping. | Overlap tests for load/unload/cancel; failure-path ownership checks; lock-order review; duplicate-load and orphan checks; I/O impact record. |
| **C4 — Exit and shutdown consistency** | Timeout, cancel, shutdown, kill, and wait races have one terminal result; unexpected helper exit is authoritative; single `Wait()` ownership holds; no double-close, send-on-closed, or cleanup-before-kill; shutdown prevents new work. | Grace and shutdown waits are bounded; no busy wait; no unnecessary file rewrite during cleanup; process and handle release is prompt. | Race tests for timeout/cancel/exit/shutdown vs active work; single-wait and double-terminal assertions; goroutine/process/handle leak checks; cleanup-ordering review; I/O impact record. |
| **C5 — Budget, backoff, and idle** | Restart ledger, backoff timers, idle deadline, and new requests have defined races; no lost wakeup, double close, timer leak, or deadlock between timer and request paths. | Timers/events replace CPU polling; restart count and backoff are bounded; no disk churn per tick; idle close occurs after the contract-defined 120 s grace and actually releases resources; the idle timer is not treated as an orphan. | Timer-vs-request overlap tests; restart/backoff bound tests; idle-grace and idle-close resource check; race/stress evidence; I/O impact record. |
| **C6 — Identity, recovery, and replacement** | Identity transitions are serialized; replacement vs pending requests has a defined winner; stale completions cannot overwrite the new generation; cap admission cannot be bypassed. | The old process is shut down, killed if necessary, and `Wait()` is observed before any new process starts. No concurrent old/new process or model residency. No redundant model-load writes; replacement is bounded. | Concurrent identity/replacement tests; no-overlap and process-count assertions; stale-generation checks; resource release review; I/O impact record. |
| **C7 — Production adapter and explicit load** | Adapter methods are safe under concurrent calls; direct cancel/unload vs synthesis is defined; cap admission is atomic; no path bypasses the limit; explicit load cannot race shutdown. | Adapter is default-off; no startup load; catalog/discovery work is cached and bounded; no unnecessary renderer/status/settings/logging/disk I/O; queues and helpers are bounded. | Adapter concurrency tests; cap-admission race tests; default-off and no-eager-load checks; resource budget and I/O impact review covering every touched application path; caller-map test. |
| **C8 — Runtime verification** | Race/stress/soak coverage for cancellation, shutdown, timeout, load/unload, retry, replacement, unexpected exit, and cap admission; lock-order, goroutine, process, handle, and temporary-file leak checks; clean repeated runs. | Measure CPU, memory, disk I/O, GPU/accelerator use where applicable, handles, and process count at idle and under load; include the declared low-core/low-memory profile; execute PR-1 through PR-4 with approved thresholds. | Pre-execution C8 plan with hardware profile, workload, warm-up, run count, duration, budgets, and thresholds; repeatable report with failure-path results; cannot be satisfied by fake-helper-only evidence where the claim is about real hardware/resources. |

## 9. Application-wide I/O evidence

### 9.1 Mandatory stage record

Every stage report must include an I/O impact record for every application path it changes or directly depends on. A stage with no such path must record an explicit N/A and have it reviewed; silence is not N/A.

Each record must state:

- path or component;
- operation (read, write, delete, scan, discover, log, or other);
- trigger and cadence;
- necessity and product/durability reason;
- bound on size, frequency, and retention;
- whether the write is model-loading, output publication, settings, cache, logging, or temporary-file handling;
- measurement or design proof;
- N/A rationale and reviewer, where applicable.

### 9.2 Application rules

- **Settings:** do not rewrite the settings file on read, status refresh, or an unchanged apply.
- **Caches:** do not rescan or rewrite a valid cache on every request; invalidate deliberately.
- **Logging:** do not emit per-frame, per-poll, or per-operation disk logs when a bounded state-change record is sufficient.
- **Status/health:** do not perform model discovery or file reads merely to answer a status request.
- **Renderer/catalog:** do not start a helper or load a model from an inspection path unless the contract and owner decision permit it.
- **Output/temporary files:** write only what is required, publish atomically where needed, and clean up deterministically.
- **Model loading:** materialization and validation writes are permitted when necessary; they must not become a per-request or per-retry rewrite.

### 9.3 Measurement

Where a stage introduces or changes I/O, the record must state how the frequency and size were measured or bounded. “No unnecessary I/O” is not demonstrated by an unmeasured assertion. C8 must include the I/O dimensions in PR-1/PR-2 where they are relevant.

### 9.4 Unrelated legacy I/O

An unrelated existing I/O problem discovered during this work is recorded as a separate finding. It does not expand the current stage and is not silently fixed outside scope. A stage may reference it but cannot use it as an excuse for adding new I/O.

## 10. Resource-budget record

Every stage must record each applicable resource in a table with these fields:

| Field | Requirement |
|---|---|
| Resource | Worker, queue, retry, timer, process, memory, buffer, file, handle, job object, storage, or log volume |
| Limit | A numeric value, or an explicit contract-provisional value with its source |
| Basis | Why this limit is safe for the workload and hardware profile |
| Measurement | Test, trace, count, or design proof used to verify it |
| Acceptance threshold | The pass/fail boundary for this stage |
| Approval | The stage reviewer/owner record that accepted the budget |

No “bounded” claim is accepted without a limit and a measurement. A limit may be provisional only where the contract labels it provisional; the source must be named. An arbitrarily high limit is not evidence of safety. The limit must not impose an arbitrary normal-operation CPU/GPU ceiling; justified concurrency and process caps remain permitted under HW-1.8.

## 11. Observability

- Bounded internal counters, state changes, and diagnostics are the default observation surface.
- Do not emit per-frame, per-poll, or per-operation disk logs; prefer bounded state changes.
- Record queue depth, rejection, retry/backoff, process count, and resource-pressure indicators where those conditions affect behavior.
- A degradation response must be one of: wait within a bound, queue within a bound, reject with a contract-permitted result, or reduce nonessential work.
- Observability must not hold lifecycle locks across I/O or allocate unbounded per-request history.
- A health/status response must not become a hidden model-discovery or disk-scan operation.
- No renderer, public API, Wails binding, or status-shape change is authorized by this section. Any external exposure requires a separate contract/owner decision.

## 12. Conflict, exceptions, and waivers

- **Contract conflict:** if a TH/HW requirement cannot be met without violating the frozen contract, stop and record the conflict. Do not silently weaken either document.
- **Stricter requirement:** a stricter TH/HW rule may be applied when compatible with the contract.
- **No safety waiver for performance:** TH-1, HW-2, synchronization, data integrity, process cleanup, lifecycle correctness, termination ordering, and contract-required limits may not be waived to improve latency, throughput, or benchmark results.
- **Budget exception:** a HW-1 budget exception is permitted only through an explicit owner decision recorded after adoption, with scope, duration, safety evidence, and a reason other than speed alone. It must not waive TH-1 or HW-2.
- **No self-approval:** the author of a change cannot approve its own exception.
- **No implied authorization:** an exception, stage acceptance, review, or this document does not authorize implementation, C8, or deployment.

## 13. Pending owner decisions

This document records engineering requirements only. It does not resolve any pending owner question.

- D10(a) and D10(b) remain pending in the owner-decision instrument.
- D1a, D9, D1b, D2, D5, and D3 remain pending.
- A conversational preference for a D10(a) result is not a recorded owner decision unless completed in the owner-decision instrument; this document neither adopts nor resolves it.
- All eight Stage 3 owner decisions remain pending.
- C1 and C2 remain unauthorized; all later stages remain unauthorized; C8 still requires separate runtime-execution authorization.

## 14. Adoption checklist

Before owner adoption, all of the following must be true:

- [ ] Independent Mimo review is complete and recorded with a verified hash.
- [ ] The review found no unresolved authority, contract, or safety contradiction.
- [ ] The owner-adoption record names this exact artifact hash and scope.
- [ ] The owning contract remains unamended.
- [ ] The owner-decision record remains unchanged and all eight decisions remain pending.
- [ ] No exception is being used to waive TH-1, HW-2, cleanup, lifecycle, or data-integrity requirements.
- [ ] The adoption record states that adoption does not authorize any stage or implementation.

Until every item is complete, this revision is **PROPOSED — NOT BINDING**. A checked box is a claim requiring evidence, not a waiver.
