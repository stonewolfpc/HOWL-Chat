# CROSS-CUTTING ENGINEERING REQUIREMENTS — Thread Safety, Hardware Safety, and Performance Preservation

**Revision:** r3 corrective  
**Date:** 2026-11-07  
**Status:** PROPOSED CORRECTIVE REVISION — NOT ADOPTED — NOT BINDING — A-1 SCOPE DECISION PENDING  
**Audited predecessor:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r2-2026-11-07.md` (SHA-256 `c9bddd10608993ae73d86caa01a7adb6dd83220f98da8dab17e46c1aa3deb8ab`, 447 lines / 36,025 bytes), preserved unchanged and audited by Mimo; r2 was not adopted.  
**Audited lineage:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md` (SHA-256 `ee38f0310042188a6821f8632e9e07bfa60dcac454a71ff2d972bb47439ac0e5`, 344 lines / 25,933 bytes), preserved unchanged; its recorded requirements remain in force as stated in §2.1.  
**Pre-adoption gates:** Independent review of this r3 revision pending; owner adoption pending; **A-1 owner scope decision pending**.  
**Governing inputs:** Exact paths and hashes are recorded in §2.2; the adoption record must cite that table.  
**Applies to:** TTS Stage 3 corrective stages **C1–C8** and the eventual implementation.  
**Authority:** Additive engineering requirement proposal. It does not amend the governing contract, the Stage 3 corrective design, the correction record, the owner-decision record, or any audit. It becomes binding only through the separate owner-adoption record in §2.

This revision responds to the read-only audit findings F1–F10, the project-owner requirement to preserve normal-operation performance, and Mimo's r2 audit findings A-1–A-11. It is not an adoption, an authorization, a stage acceptance, or a decision on any pending owner question. It does not select an outcome for A-1.

## 1. Corrective basis

| Audit finding | Resolution in this revision |
|---|---|
| **F1 — authority** | §2 defines a separate, hash-bound owner-adoption mechanism and records the verified input hashes. It also states that the original recorded requirements remain in force until r3 adoption. |
| **F2 — C1 handshake contradiction** | C1 now names the only start triggers: explicit TTS/synthesis, catalog, or status requests under contract §17.7.2. A start from direct `unload`, direct `Cancel`, unknown-capability inspection, startup, or health is not authorized; completing a handshake is not a start trigger. See §8, C1. |
| **F3 — C6 overlap** | Replacement requires shutdown, termination if necessary, and observed `Wait()` before any new process starts; no concurrent old/new residency. See §8, C6. |
| **F4 — timeout/helper-exit races** | Timeout and unexpected helper exit are first-class TH-1.2 events and C4/C8 evidence. See §4.2, §8 C4, and §8 C8. |
| **F5 — 120 s idle grace** | The idle timer and helper may remain through the contract-defined 120 s grace; premature shutdown is not required. See §5.3 and §8 C5. |
| **F6 — app-wide I/O evidence** | §9 defines a mandatory I/O impact record for every stage and every touched application path. See AC-7 and §8. |
| **F7 — measurable budgets** | §10 defines a resource-budget record; “bounded” requires a numeric value or explicit contract-provisional value, a measurement method, and an acceptance threshold. |
| **F8 — exception loophole** | §12 prohibits performance waivers of TH-1, HW-2, contract cleanup/lifecycle/data-integrity requirements, and process ordering. Any HW-1 budget exception requires a recorded owner decision with scope, expiry, and safety evidence. |
| **F9 — observability surface** | §11 makes bounded internal diagnostics the default. No renderer, public API, or status-shape change is authorized without separate contract/owner authority. |
| **F10 — C8 pass/fail** | C8 retains the frozen design's pre-execution evidence matrix, including hardware/workload details where the design requires them. PR-1–PR-4 are not C8 acceptance conditions; A-1 remains an open owner scope decision. See §6.4–§6.5. |
| **New requirement — performance preservation** | HW-1.8 retains normal-operation performance preservation. HW-2.4 defines the PR-1–PR-4 measurement framework, but assigns it no acceptance effect until A-1 is decided. |

No finding justifies abandoning the safety requirements. The predecessor remains the audited revision; this file is the proposed corrective revision for independent review.

### 1.1 Mimo r2 audit findings A-1–A-11

Mimo's r2 audit is review evidence, not authority. Its full artifact is `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r2-2026-11-07.md` (SHA-256 `567bc86c3d8b2a78b9687131d5937d77a9736393bed6a383dd0aba5ff948994b`, 221 lines / 29,130 bytes).

| Finding | Severity | r3 disposition |
|---|---|---|
| **A-1 — C8 performance-gate conflict** | Blocking | The PR-1–PR-4 C8 mandate is removed. A-1 is recorded as an open owner scope decision in §6.5; r3 does not choose between a separate performance gate and a formal C8 amendment. |
| **A-2 — all eight decisions required pending** | Major | §13 and §14 no longer make “all eight remain pending” an adoption condition. Adoption neither requires nor resolves any owner decision. |
| **A-3 — predecessor status unresolved** | Major | §2.1 states that the original recorded requirements remain in force until r3 adoption and that r3 supersedes them only where the texts differ after adoption. |
| **A-4 — no governing-input hashes** | Major | §2.2 records the exact r3 inputs and hashes; the adoption record must cite that table. |
| **A-5 — handshake trigger unnamed** | Blocking for C1 | §8 C1 names the permitted request triggers and explicitly excludes direct `unload`, direct `Cancel`, unknown-capability inspection, startup, and health. |
| **A-6 — inconsistent PR timing** | Major | §6.4 assigns PR-1–PR-4 no stage or C8 acceptance effect in r3. Stages record an AC-8 impact statement only. |
| **A-7 — scope narrowed** | Minor | §3.1 restores “relies on … or is directly integrated with” and retains the anti-refactor guard. |
| **A-8 — exception can fail PR-3** | Minor | §12 states that no HW-1 exception may create a condition that fails a HW-2 pass condition; the normal-operation throttle prohibition is treated as part of HW-2.1 for that purpose. |
| **A-9 — PR-3 has no N/A or signal** | Minor | §6.4 adds a reviewed N/A rule for stages without adaptive control and requires a named contention signal and detection latency when controls exist. |
| **A-10 — AC-7 N/A over-breadth** | Minor | AC-7 now requires N/A only for each in-scope path category not touched, by reference to §9.1. |
| **A-11 — circular approval field** | Minor | §10 names the independent reviewer as budget approver at the stage gate, and the owner only for exceptions; C8 performance thresholds, if separately authorized, are approved before execution. |

## 2. Authority and owner adoption

### 2.1 Current status

This document is **PROPOSED — NOT BINDING**. Until the owner completes a separate adoption record, it may be used as a review checklist but cannot be used to accept or reject a stage, authorize implementation, or waive any contract requirement.

The original recorded requirements (`REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md`) remain in force, unamended, until this revision is adopted. Nothing in this revision suspends the predecessor's acceptance rule in the interim. r2 was audited but was not adopted and has no operative effect. Adoption of r3 supersedes the original only where the two texts differ; where r3 is silent, the original remains in force.

The governing contract remains controlling. This document does not amend §17 or any adopted clarification. If a requirement conflicts with the contract, the conflict must be recorded and escalated; neither document may be silently overridden.

### 2.2 Separate adoption record required

A document cannot adopt itself, and inserting this file’s final hash into this file would change the hash. Adoption must therefore be recorded separately and must bind the exact artifact. The adoption record must contain:

- owner identity or accountable role;
- adoption date;
- exact path and SHA-256 of this file;
- adopted scope, including C1–C8 and the eventual implementation;
- the exact governing-input table below, with hashes rechecked at adoption;
- confirmation that the governing contract and frozen design inputs match the table and were not amended;
- the independent review record for r3 and its verified hash;
- the A-1 owner scope decision, if any performance gate is to have acceptance effect; if A-1 is unresolved, the adoption record must state that PR-1–PR-4 are inactive;
- any approved exceptions, with scope, duration, and safety evidence;
- an explicit statement that adoption does not authorize any stage or implementation.

Verified inputs at the time of drafting:

| Input | Path | SHA-256 | Lines | Bytes |
|---|---|---|---:|---:|
| Original recorded requirements | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md` | `ee38f0310042188a6821f8632e9e07bfa60dcac454a71ff2d972bb47439ac0e5` | 344 | 25,933 |
| Audited r2 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r2-2026-11-07.md` | `c9bddd10608993ae73d86caa01a7adb6dd83220f98da8dab17e46c1aa3deb8ab` | 447 | 36,025 |
| Mimo r2 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r2-2026-11-07.md` | `567bc86c3d8b2a78b9687131d5937d77a9736393bed6a383dd0aba5ff948994b` | 221 | 29,130 |
| Governing contract | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 | 182,904 |
| Corrective proposal | `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 | 69,168 |
| Correction record | `docs/review/CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` | `8fe9356cc304f3135657d19a7f2ea6badd9cf2233585a03bd0bd5c4f6355e076` | 172 | 16,190 |
| Owner-decision instrument | `docs/review/OWNER-DECISIONS-howl-tts-stage3-2026-11-07.md` | `24d3344554e7ed6a465c8b7bdbca4461cda23f78353aa7626e3cf4ba7df122ba` | 610 | 45,891 |

A hash mismatch in any row reopens review. The table is not a substitute for the separate r3 hash: the adoption record must carry both.

### 2.3 Adoption block — left blank

**Owner:**  
**Adoption date:**  
**Artifact path:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r3-2026-11-07.md`  
**Artifact SHA-256:**  
**Adopted scope:**  
**Governing-input table (§2.2) rechecked:** [ ] Yes  [ ] No  
**Independent r3 review record / hash:**  
**A-1 owner decision (if PR gate has acceptance effect):**  
**Contract-amendment confirmation:**  
**Approved exceptions:**  
**Decision:** [ ] Adopt  [ ] Reject  [ ] Return for correction

### 2.4 Pre-adoption rule

Owner adoption must not occur before the independent review of r3 is complete and recorded. A review finding that identifies a safety, authority, or contract contradiction reopens this revision rather than being treated as an adoption prerequisite waived by schedule. Mimo's r2 audit is complete; it is not the independent review of r3.

A-1 must be decided before any PR-1–PR-4 result is given acceptance effect. If A-1 remains open, r3 may be considered only as a non-performance requirements revision, with PR-1–PR-4 expressly inactive; it cannot be adopted as a performance-gating instrument.

No adoption, review, or stage acceptance authorizes C1, C2, C3, C4, C5, C6, C7, C8, implementation, runtime execution, or benchmarking.

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

It also applies to code changed by those stages and code those stages rely on for shared state, process ownership, scheduling, cleanup, or user-visible status, including code directly integrated with those paths.

“All code those stages rely on” is limited to code the stage changes, relies on, or directly integrates. It does not authorize an audit or retrofit of unrelated legacy code, third-party internals, or the entire application. Unrelated I/O debt discovered during the work is recorded separately under §9.4.

### 3.2 Application-wide I/O scope

HW-1 applies to every new or changed application path, not only the TTS helper. This includes settings, caches, model/asset preparation, status and health observation, renderer and picker paths, logging and diagnostics, helper lifecycle, and temporary/output files.

“Application-wide” means all paths the implementation changes, relies on, or directly integrates with. It does not mean that every pre-existing unrelated write in the application must be fixed by this program. The distinction must be recorded per stage in the I/O impact record (§9).

### 3.3 What this document does not do

This revision does not:

- amend the frozen contract or its §17 rules;
- amend the Stage 3 corrective design or correction record;
- amend the owner-decision record;
- resolve D10(a), D10(b), D1a, D9, D1b, D2, D5, or D3;
- decide A-1, either by creating a separate performance gate or by amending C8's scope;
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
- Adaptive controls, where present, must activate in response to actual, measured contention, must be reversible, and must deactivate when contention clears. A stage with no adaptive control records that fact in its AC-8 impact statement; if a performance gate is later authorized, PR-3 is N/A with review under §6.4.
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

### 6.4 HW-2.4 Performance-regression framework — acceptance assignment pending A-1

PR-1–PR-4 define a candidate performance-preservation measurement framework. They are not acceptance conditions for C8 or for any stage in this revision. The frozen design's C8 exclusion remains controlling: C8 does not make a performance claim, and no PR result may be used to accept or reject C8. The A-1 owner scope question is unresolved; §6.5 records the two options without selecting one.

The framework remains available for a separately authorized performance-preservation gate if the owner selects that option. A stage that changes scheduling, concurrency, throttling, caching, or normal-operation work records an AC-8 impact statement naming the PR dimensions its change affects. It does not run the full PR suite for stage acceptance. A stage may describe a PR subset in its AC-8 impact statement, but it does not execute the suite or use results as acceptance evidence while A-1 is open.

| Check | Measurement | Pass condition |
|---|---|---|
| **PR-1 — Normal-operation overhead** | No injected contention; baseline versus candidate; latency (median and p95), throughput, CPU, memory, disk I/O, handles, and process count where applicable. GPU/accelerator dimensions apply only if an authorized native-provider scope includes them; otherwise record N/A with review. | Candidate does not exceed the approved regression threshold; no arbitrary CPU/GPU ceiling or unnecessary throttle is active |
| **PR-2 — Contention degradation** | Same workload under a declared contention profile; baseline versus candidate | Degradation is bounded by the approved threshold, and safety behavior remains correct |
| **PR-3 — Adaptive control** | Activation and deactivation under measured contention | Controls activate only after actual contention, are reversible, and deactivate when contention clears; no permanent normal-operation throttle. A stage with no adaptive control records PR-3 as N/A with review. Where controls exist, the contention signal must be named — queue depth, rejection count, or measured latency against the declared profile — and the detection latency must be stated so “after actual contention” is checkable. |
| **PR-4 — Safety regression** | Race, timeout, cancellation, exit, cleanup, and resource-limit checks under the same workload | Any safety regression blocks acceptance regardless of PR-1/PR-2 results |

If a separate performance gate is later authorized, its numeric regression thresholds, contention profile, run count, and warm-up period must be recorded and approved by the independent reviewer before execution. This document does not invent them. A report that presents measurements without a pass/fail comparison is not performance-gate evidence. No PR result may override PR-4 or any safety requirement.

### 6.5 A-1 open owner scope decision — performance gate

**Status: OPEN — OWNER DECISION REQUIRED. This section does not decide A-1.**

Mimo's r2 audit A-1 found that r2 made PR-1–PR-4 mandatory C8 acceptance evidence while the frozen design explicitly excludes performance claims from C8. r3 removes that conflict by assigning no acceptance effect to PR-1–PR-4. It does not choose the owner-level scope question.

**Option A — separate performance-preservation gate (recommended for investigation, not selected).** PR-1–PR-4 run only under a separate, explicitly authorized gate with its own plan, independent approval, thresholds, and acceptance record. C8 remains unchanged and excludes performance claims. The performance gate does not gate C8 contract closure. This preserves the design boundary and gives the performance requirements a meaningful home.

**Option B — formal C8 amendment.** The owner explicitly amends the design's C8 exclusion, records that amendment in the adoption record, and makes PR-1–PR-4 C8 acceptance conditions. This cannot be described as a supplement and cannot be paired with an unqualified “frozen inputs were not amended” certification.

No third option is available: r3 must not leave the C8 exclusion unmentioned while giving PR evidence C8 effect.

Until the owner selects A or B, PR-1–PR-4 are inactive. They may be reviewed as a measurement design, but no PR result may accept, reject, block, or delay C8 or any stage. No performance measurement or benchmark is authorized.

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
| **AC-7 — Application I/O evidence** | Every stage records its I/O impact under §9, including explicit N/A for each in-scope path category it did not touch, per §9.1. |
| **AC-8 — Performance preservation** | The stage records whether its changes affect normal-operation performance, naming the PR dimensions affected. C8 does not apply PR-1–PR-4; the performance gate remains unassigned pending A-1. |
| **AC-9 — Applicability review** | A TH-1 or HW-1 “not applicable” determination requires an explicit reviewed entry; blanket N/A is not accepted. |
| **AC-10 — No stage substitution** | C8 may aggregate evidence but cannot substitute for missing C1–C7 evidence. |

A stage that has functional tests but no applicable TH/HW evidence is not accepted. A stage with no declared resource budget is not accepted.

## 8. C1–C8 acceptance overlay

The following table supplements—and does not replace—each stage's frozen functional acceptance criteria and every pending owner-decision gate. The D10(a)/D10(b), D1a, D9, D1b, D2, D5, and D3 gates remain in force. It adds no performance acceptance condition: C8 remains the frozen design's evidence matrix, and PR-1–PR-4 are inactive pending A-1.

| Stage | TH-1 acceptance focus | HW-1 / HW-2 / performance focus | Required evidence |
|---|---|---|---|
| **C1 — Capability negotiation** | Capability snapshot publication is atomic and single-writer; concurrent negotiation, gating, and reads do not observe torn or stale fields; handshake completion and failure have one terminal outcome. | A start is permitted only for an explicit TTS/synthesis, catalog, or status request authorized by contract §17.7.2 (including the `ProviderStatus` and `ProviderCatalog` rows); capability knowledge obtained by that request's handshake is incidental. No start is authorized from direct `unload`, direct `Cancel`, an optional-capability inspection with unknown state, startup, or health, and no start is authorized merely to complete a handshake or discover a capability. No model load or disk write from negotiation or inspection. Parsing and descriptor storage are bounded. | Trigger-to-start mapping for TTS/synthesis, catalog, and status; structural no-frame/no-spawn tests for direct `unload`, direct `Cancel`, and unknown-capability inspection; authorized-trigger start evidence; no-model-load/no-write evidence; allocation/limit review; I/O impact record (likely N/A). A test that starts a helper from `unload` or direct `Cancel` fails C1. |
| **C2 — Status projection** | Multi-field status is a coherent snapshot; lifecycle transitions cannot produce mixed old/new fields; status reads do not hold locks across I/O or block transitions indefinitely. | Status is bounded-cost and performs no model discovery, file read, or high-frequency log write; buffers and counters are bounded; observability is internal §11 diagnostics only. | Concurrent status-vs-transition tests; no-torn-snapshot assertions; blocking-call review; I/O impact record; evidence that no public status/renderer/API field was added. |
| **C3 — Lifecycle and load seam** | Load/unload/cancel/synthesis have defined overlap and winner rules; one active load; lock order between adapter, client, and lifecycle state is documented; no orphaned worker on failure. | No eager spawn or model load from inspection/status/startup; no duplicate load; waits are bounded; no write cadence added by lifecycle bookkeeping. | Overlap tests for load/unload/cancel; failure-path ownership checks; lock-order review; duplicate-load and orphan checks; I/O impact record. |
| **C4 — Exit and shutdown consistency** | Timeout, cancel, shutdown, kill, and wait races have one terminal result; unexpected helper exit is authoritative; single `Wait()` ownership holds; no double-close, send-on-closed, or cleanup-before-kill; shutdown prevents new work. | Grace and shutdown waits are bounded; no busy wait; no unnecessary file rewrite during cleanup; process and handle release is prompt. | Race tests for timeout/cancel/exit/shutdown vs active work; single-wait and double-terminal assertions; goroutine/process/handle leak checks; cleanup-ordering review; I/O impact record. |
| **C5 — Budget, backoff, and idle** | Restart ledger, backoff timers, idle deadline, and new requests have defined races; no lost wakeup, double close, timer leak, or deadlock between timer and request paths. | Timers/events replace CPU polling; restart count and backoff are bounded; no disk churn per tick; idle close occurs after the contract-defined 120 s grace and actually releases resources; the idle timer is not treated as an orphan. | Timer-vs-request overlap tests; restart/backoff bound tests; idle-grace and idle-close resource check; race/stress evidence; I/O impact record. |
| **C6 — Identity, recovery, and replacement** | Identity transitions are serialized; replacement vs pending requests has a defined winner; stale completions cannot overwrite the new generation; cap admission cannot be bypassed. | The old process is shut down, killed if necessary, and `Wait()` is observed before any new process starts. No concurrent old/new process or model residency. No redundant model-load writes; replacement is bounded. | Concurrent identity/replacement tests; no-overlap and process-count assertions; stale-generation checks; resource release review; I/O impact record. |
| **C7 — Production adapter and explicit load** | Adapter methods are safe under concurrent calls; direct cancel/unload vs synthesis is defined; cap admission is atomic; no path bypasses the limit; explicit load cannot race shutdown. | Adapter is default-off; no startup load; catalog/discovery work is cached and bounded; no unnecessary renderer/status/settings/logging/disk I/O; queues and helpers are bounded. | Adapter concurrency tests; cap-admission race tests; default-off and no-eager-load checks; resource budget and I/O impact review covering every touched application path; caller-map test. |
| **C8 — Runtime verification** | Race/stress/soak coverage for cancellation, shutdown, timeout, load/unload, retry, replacement, unexpected exit, and cap admission; lock-order, goroutine, process, handle, and temporary-file leak checks; clean repeated runs. | C8 uses the frozen design's evidence matrix. No performance claim, throughput/latency claim, GPU/accelerator claim, or PR-1–PR-4 condition is added. Resource, handle, and process observations are limited to the safety evidence the design already requires. | Pre-execution C8 evidence plan; commands, exit codes, logs, and hashes; every row names its environment and process boundary; Windows-specific rows close only through a Windows real-child run; no performance claim; PR results are not C8 evidence. |

## 9. Application-wide I/O evidence

### 9.1 Mandatory stage record

Every stage report must include an I/O impact record for every application path it changes, relies on, or directly integrates with. A stage with no such path must record an explicit N/A and have it reviewed; silence is not N/A.

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

Where a stage introduces or changes I/O, the record must state how the frequency and size were measured or bounded. “No unnecessary I/O” is not demonstrated by an unmeasured assertion. If a separate performance gate is later authorized, it must include the I/O dimensions in PR-1/PR-2 where they are relevant. C8 itself makes no performance claim and runs no PR check.

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
| Approval | The independent reviewer who approves the budget at the stage acceptance gate; the owner only for an exception under §12. The change author cannot approve. |

No “bounded” claim is accepted without a limit and a measurement. A limit may be provisional only where the contract labels it provisional; the source must be named. An arbitrarily high limit is not evidence of safety. The limit must not impose an arbitrary normal-operation CPU/GPU ceiling; justified concurrency and process caps remain permitted under HW-1.8.

Budget approval is recorded at the stage acceptance gate by the independent reviewer, after the budget is prepared but before acceptance; the same approval does not create its own input. The owner is the approver only for an exception under §12. The change author may not approve. If a separate performance gate is later authorized, its thresholds are approved by the same independent reviewer before execution, not at C8 acceptance.

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
- **No exception may fail HW-2:** an HW-1 exception may not create a condition that fails a HW-2 pass condition. For this purpose, HW-1.8's prohibition on normal-operation throttling is treated as part of HW-2.1. An exception that would require a permanent normal-operation throttle is not available.
- **No self-approval:** the author of a change cannot approve its own exception.
- **No implied authorization:** an exception, stage acceptance, review, or this document does not authorize implementation, C8, or deployment.

## 13. Pending owner decisions

This document records engineering requirements only. It resolves no owner question.

- At drafting, all eight Stage 3 owner decisions are pending: D10(a), D10(b), D1a, D9, D1b, D2, D5, and D3.
- This is a status report, not an adoption condition. Adoption neither requires any decision to remain pending nor treats an answered decision as authority for this revision.
- If any decision is answered before adoption, r3 does not incorporate that answer; the stage gates remain as stated in the owner-decision instrument.
- A conversational preference for a D10(a) result is not a recorded owner decision unless completed in the owner-decision instrument; this document neither adopts nor resolves it.
- A-1 is a separate open owner scope decision. It is not one of the eight Stage 3 decisions and is not resolved here.
- C1 and C2 remain unauthorized; all later stages remain unauthorized; C8 still requires separate runtime-execution authorization.
- No PR-1–PR-4 result has acceptance effect while A-1 is open.

## 14. Adoption checklist

Before owner adoption, all of the following must be true:

- [ ] Independent review of r3 is complete and recorded with a verified hash.
- [ ] The review found no unresolved authority, contract, or safety contradiction.
- [ ] The owner-adoption record names this exact artifact hash and scope.
- [ ] The §2.2 governing-input hashes were rechecked and match.
- [ ] The owning contract remains unamended.
- [ ] This revision resolves no owner decision; the adoption record does not require or imply that any decision remains pending, and does not cite an answered decision as authority.
- [ ] A-1 is explicitly decided before any PR-1–PR-4 result is given acceptance effect; if A-1 is open, the adoption record states that PR-1–PR-4 are inactive.
- [ ] No exception is being used to waive TH-1, HW-2, cleanup, lifecycle, or data-integrity requirements.
- [ ] The adoption record states that adoption does not authorize any stage or implementation.

Until every item is complete, this revision is **PROPOSED — NOT BINDING**. A checked box is a claim requiring evidence, not a waiver.
