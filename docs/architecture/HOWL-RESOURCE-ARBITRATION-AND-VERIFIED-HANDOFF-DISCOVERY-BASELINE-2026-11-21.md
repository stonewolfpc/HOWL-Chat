# Howl — Resource Arbitration and Verified Handoff Discovery Baseline

## Discovery baseline `HRA-DB-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Architectural discovery baseline; requirements and unresolved questions only |
| Starting scope | `HFS-AUDIT-2026-11-21-01` |
| Readiness | `HFR-CP-2026-11-21-01` |
| Program map | `HCDM-2026-11-21-01` |
| Phase 1 | `HF1-PLAN-2026-11-21-01` — still proposed, not authorized |
| Prior architecture | `AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md`; `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` |
| Production/code/schema/migration changes | **0** |
| Tests/models/benchmarks/process runs | **0** |
| Implementation authorization | **None** |
| Status | **Discovery baseline.** Owner-approved requirements, architectural proposals, unresolved questions and rejected alternatives are labeled separately |

**Labels:** **Owner-approved** · **Observed** · **Source-derived** · **Proposed** · **Unresolved** · **Rejected** · **Calculated** · **Unestablished**.

---

## 1. Decision

**HRA-01 (Owner-approved requirement).** Howl must treat **cross-modality resource arbitration and verified model handoff** as a foundational application capability, not as a model-manager convenience and not as an optional optimization.

The central principle is:

> **Selection is intent. Admission is a decision. Loading is an allocation. Ready is an evidence claim. Unloading is not release.**

Each responsibility belongs to a different boundary. None may claim another boundary's work is complete.

**HRA-02 (Proposed).** The handoff protocol belongs to the application's resource-governance architecture. The model lifecycle adapter executes drain, unload, load and probe mechanics; it does not decide global admission. Chat requests a capability; it does not own hardware policy or runtime internals.

**HRA-03 (Owner-approved requirement).** The architecture must accommodate future audio, video, vision, embedding and other model-consuming workloads without redesigning the handoff protocol for each modality.

---

## 2. Corrected modality scope ledger

**HRA-04 (Owner-approved — current scope).** The following table supersedes older statements in `HCDM-2026-11-21-01` that described image generation as deferred and audio generation as removed. Those statements are **Superseded** for current scope.

| Capability | Current scope | Planning consequence |
|---|---|---|
| **Audio capture/input** | Separate capability; current behavior to be assessed | Must have its own pipeline, lifecycle, resource envelope and failure modes. It is not audio generation and not TTS. |
| **Audio generation** | **Complete** | No remaining delivery work is claimed here. Completion is an owner statement; the evidence record must be linked separately. |
| **TTS** | **Deferred, but required for chat finalization** | It stays out of current scheduling and out of M0A/M0B. Chat cannot be declared final while a required TTS path is absent. Re-entry needs a dated owner decision and its own work package. |
| **Image generation** | **Nearly complete; remaining work stays in scope** | T0.3/T0.4/T2.3/T3.9/T3.10 are no longer classified as deferred. Remaining work needs an evidence-based closure list, not a blanket “deferred” label. |
| **Vision** | **Largely complete; remaining gaps need assessment** | Vision input's remaining gaps must be inventoried before it can be called complete. Package presence is not acceptance. |
| **Video** | **Incomplete; remains in scope** | It needs its own capability contract, workspace/storage estimate and lifecycle boundary. It is not a variant of still-image generation for resource purposes. |
| **2D/3D visual-character rendering** | **Removed from Howl scope** | No visual-character workstream, dependency, resource reservation or acceptance gate. Historical records remain historical. |

**HRA-05 (Calculated — separation rule).** Audio capture/input, audio generation and TTS are three separate entries because they have different pipelines, latency requirements, model lifecycles and failure modes. Collapsing them would make resource accounting and readiness claims dishonest.

**HRA-06 (Owner-approved requirement).** An older ledger must not silently override these decisions. Any document that still says image generation is deferred, audio generation is removed, or TTS is optional must be treated as superseded for current scope.

**HRA-07 (Calculated — scope is not acceptance).** “Complete” and “nearly complete” are owner scope statements. They do not establish E3 evidence. The readiness levels in `HFR-CP-2026-11-21-01` still apply.

---

## 3. Why this is foundational

**HRA-08 (Owner-approved requirement).** Chat is not only a text-generation interface. It is the orchestration point for capabilities that consume the same physical machine, have different lifecycles, and can compete for memory, compute, storage and thermal headroom.

**HRA-09 (Proposed).** The system needs a hardware-aware resource and safety authority that can:

- observe what the machine is doing;
- estimate what the next operation needs;
- reserve capacity before allocation;
- decide whether to coexist, queue, transition, warn or refuse;
- drive a safe drain and handoff;
- verify release and readiness as far as the platform permits;
- report uncertainty honestly;
- preserve user work across transitions;
- retain evidence for every consequential decision.

**HRA-10 (Owner-approved requirement).** A model fitting into memory does not mean Howl should load it there. “Fits” is a prediction; “safe and sustainable” is the admission standard.

**HRA-11 (Calculated — correction to an earlier draft).** Howl must **not** impose an arbitrary one-model-at-a-time rule. `AMENDMENT-01` already rejects that. Multiple resident models are allowed when resources and backend behavior safely support them. Conservative serialization is a policy for constrained or uncertain conditions, not a permanent architectural limit.

---

## 4. Eight boundaries

**HRA-12 (Proposed — conceptual decomposition, not a final component diagram).**

| Boundary | Owns | Must not own |
|---|---|---|
| **1. Resource and safety governance** | Hardware truth, estimates, admission, reservations, headroom, priority, safety policy, cross-consumer arbitration | Model semantics, chat history, modality internals |
| **2. Model lifecycle and runtime adapters** | Artifacts, runtime processes, drain, unload, load, warmup, readiness probes, release evidence | Global admission or permission to declare the machine safe |
| **3. Conversation orchestration** | Turn state, context assembly, capability requests, ordering, cancellation, preservation of user work | Direct hardware control, modality implementation, memory storage |
| **4. Modality services** | Vision, audio input, audio generation, TTS, image generation, video, embeddings/rerank pipelines and output contracts | Conversation history or global resource policy |
| **5. Memory and knowledge** | MEMORY-09, legacy lorebook, retrieval, indexes, embeddings, provenance, deletion | Chat turn state or model loading |
| **6. State, persistence and recovery** | Canonical conversation state, event history, snapshots, idempotency, crash recovery, multi-instance file ownership | Model choice or resource policy |
| **7. User control and transparency** | Status, consequences, manual strategy, confirmation, degradation notices, recovery actions | Hidden automatic decisions |
| **8. Evidence and verification** | Telemetry, decision logs, release proof, readiness proof, failure records, test evidence | Converting uncertainty into a success claim |

**HRA-13 (Proposed — responsibility rule).** A governor may approve a reservation without knowing how a runtime releases allocations. A runtime adapter may report that its process exited without having authority to declare the whole machine safe for another model. Chat may request a capability without owning its implementation. Evidence may record a claim without granting permission to act.

---

## 5. Relationship to the existing architecture package

**HRA-14 (Observed).** The existing architecture is not empty:

- `AMENDMENT-01` already establishes unified hardware intelligence, model selection versus residency, leases, release verification, no arbitrary one-model limit, multi-GPU as a first-class requirement, and manual override rules.
- `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` already defines a hardware inventory, telemetry confidence, the 10% reserve, a lease authority, lifecycle states, release verification and acceptance design.
- `HCDM-2026-11-21-01` already contains a model-lifecycle milestone and a contract register.
- `HF1-PLAN-2026-11-21-01` correctly excludes hardware lease implementation from Phase 1.

**HRA-15 (Calculated — this baseline's contribution).** This baseline extends that package with:

1. an explicit cross-modality handoff protocol, rather than lifecycle mechanics alone;
2. a corrected current modality ledger;
3. logical reservation versus physical allocation as a first-class distinction;
4. cross-instance reservation coordination as an unresolved architecture problem;
5. a generalized evidence ladder for all state transitions;
6. a traceability chain from original requirement to implementation evidence.

It does not replace `AMENDMENT-01` or the hardware design. Where they conflict, the later owner decision and this baseline control scope; the earlier design remains the detailed architectural proposal.

---

## 6. Handoff protocol

**HRA-16 (Proposed — required phases).** A handoff is a sequence of evidence-bearing transitions, not a single switch call.

| Phase | Question | Minimum evidence | Failure behavior |
|---|---|---|---|
| **Intent** | What capability is required, by which owner, for which work? | Capability request, priority, user-work snapshot reference | Reject invalid request; preserve work |
| **Assess** | What does the incoming capability require? | Weights, context/KV, runtime, workspace, storage, thermal estimate, confidence | Mark estimate unknown; apply uncertainty policy |
| **Admit** | May it coexist, must something transition, should it defer, or must it be refused? | Current allocations, reservations, headroom, policy decision | Queue, ask, refuse, or re-plan; never double-spend |
| **Reserve** | Has capacity been claimed before allocation? | Lease/reservation with owner, expiry and idempotency key | Do not allocate without a valid reservation |
| **Drain** | Is the outgoing capability quiescent? | No new work admitted; in-flight work completed or explicitly interrupted; state flushed | Keep old capability resident; report busy/refused |
| **Release** | Did the outgoing runtime release what it owned? | Process exit, job/handle closure, cache/workspace cleanup, memory observation | Mark partial or unknown; retain ledger |
| **Verify release** | Is release proven or unverified? | Measured drop, process/job evidence, or explicit unknown | Apply safe policy; do not claim freed capacity |
| **Load** | Is the incoming artifact valid and initializable? | Integrity check, runtime start, configuration applied | Explicit failed state; preserve selection; no silent substitution |
| **Warm/probe** | Is it ready for this workload? | Workload-specific readiness probe | Do not dispatch; report not-ready |
| **Commit** | May work proceed? | Ready state plus valid reservation | Release or downgrade reservation depending on policy |
| **Recover** | What happens after failure? | Named state, cleanup result, preserved work, retry/restore/refuse decision | No phantom ready; no destructive loop |

**HRA-17 (Owner-approved requirement).** Selection must not immediately allocate hardware. A selected model may be queued or unavailable; its selection survives without residency.

**HRA-18 (Owner-approved requirement).** Unload means actual release as far as the platform can demonstrate. Stopping inference is insufficient if weights, runtime allocations, caches, workspaces or job-owned processes remain resident.

**HRA-19 (Owner-approved requirement).** Verify before loading. Howl must check that the outgoing capability released what it was supposed to release, or explicitly report that verification is unavailable.

**HRA-20 (Proposed — drain is not optional).** Drain is a distinct phase: stop admitting new work, finish or explicitly interrupt the active operation, flush durable state, snapshot user work, release context/caches, then unload. A busy model is not “safe to kill.”

**HRA-21 (Proposed — no blind handoff).** If telemetry or release verification is unavailable, Howl must report uncertainty and follow a defined safe policy. It must not convert unknown into free, and must not claim a verified handoff.

**HRA-22 (Owner-approved requirement).** Failed handoffs recover cleanly. If the incoming model fails to load or become ready, Howl must not leave a misleading ready state, must preserve user work, and must not silently substitute another model or placement.

---

## 7. Reservation model: logical versus physical

**HRA-23 (Proposed — core distinction).** A **logical reservation** is Howl's accounting commitment. A **physical allocation** is what the operating system, driver, runtime and other applications actually do. A reservation does not guarantee that the hardware will honor the estimate.

**HRA-24 (Proposed).** The resource authority owns reservations. The model lifecycle adapter owns runtime release/load mechanics. The evidence layer records what was observed. No component may treat its own claim as proof from another component.

### 7.1 Reservation questions

**HRA-25 (Unresolved — `DQ-HRA-01`). Which component owns a reservation?**

Proposed direction: the resource authority owns the lease record; the requesting workload owns the obligation to use and release it; the lifecycle adapter owns execution against it. The authority is the only writer of admission state.

**HRA-26 (Unresolved — `DQ-HRA-02`). Can multiple operations share one reservation?**

Proposed direction: only through an explicit parent/child budget. No implicit sharing across unrelated operations. If sharing is allowed, one child exceeding its share must not silently consume another child's reserved capacity. The final policy needs a decision.

**HRA-27 (Unresolved — `DQ-HRA-03`). What happens when an operation exceeds its estimate?**

Proposed direction: monitor actual usage against the reservation; if it exceeds the declared envelope, pause/re-plan/degrade/abort according to policy and surface the overrun. Never silently continue while claiming the original reservation is valid.

**HRA-28 (Unresolved — `DQ-HRA-04`). Can a reservation survive a crash?**

Proposed direction: reservations are durable enough to reconcile, but a crash makes them untrusted until checked. On restart, reconcile the durable lease record against process/job ownership and a fresh hardware snapshot. Expired leases may be reclaimed only when the owner is known gone or the evidence supports it. A stale reservation must not outvote the device.

**HRA-29 (Unresolved — `DQ-HRA-05`). How are orphaned reservations detected?**

Proposed direction: lease heartbeat/expiry, process/job ownership, registered allocation reconciliation, and telemetry. Orphans are reported and recovered through a controlled path. They are not killed merely because they are unexplained.

**HRA-30 (Unresolved — `DQ-HRA-06`). What happens when another application consumes memory unexpectedly?**

Proposed direction: existing non-Howl usage is part of the observed baseline. Re-check at commit; re-plan or refuse if headroom is no longer safe. Howl cannot guarantee capacity that another process can take away. The user must see the contention when it matters.

**HRA-31 (Unresolved — `DQ-HRA-07`). How do two Howl instances coordinate without claiming the same capacity?**

This is the hardest open question. A per-process governor is not sufficient: two instances can each hold a logically valid ledger and both believe the same free bytes are available. T0.5 process/job ownership prevents one instance from killing the other's child; it does **not** prevent double-spending the machine's memory.

Candidate directions, none selected:

| Option | Strength | Cost/risk |
|---|---|---|
| Single supported instance; second instance refuses heavy work or runs read-only | Strong safety | Limits a use case the product may want |
| Shared local broker with one lease authority | Strong accounting | New process/lifecycle, failure and startup complexity |
| File/socket lease registry with atomic claims | No permanent broker | Crash recovery, stale locks and portability |
| Per-instance observation with conservative partitioning | Simpler | Cannot guarantee against races; must be honest about uncertainty |
| OS-level ownership plus per-instance caps | Uses platform guarantees | Does not create a global reservation ledger by itself |

**HRA-32 (Calculated).** Cross-instance arbitration is a first-class design decision, not an implementation detail. It must be resolved before any claim that concurrent instances are resource-safe.

**HRA-33 (Proposed — physical evidence wins).** A logical ledger may say capacity is free; observed physical usage may contradict it. Howl must not overrule an unsafe observed state with its own accounting. When observation is unavailable, it must not claim safety from the ledger alone.

### 7.2 Proposed concurrency rules

**HRA-33A (Proposed — not yet approved).**

| Rule | Requirement |
|---|---|
| **Single admission writer** | One authority writes admission/reservation state; consumers do not create parallel ledgers. |
| **Atomic claim** | A reservation is recorded atomically against one coherent snapshot; the same free bytes cannot be promised twice. |
| **Residency ≠ execution** | Multiple models may be resident while execution is serialized, if resources support residency and policy permits it. Neither state implies the other. |
| **Busy is not killable** | A busy capability is drained or refused under policy. User-visible work is not silently destroyed to satisfy a new request. |
| **Priority with fairness** | Interactive work may outrank background embeddings/reranking, but a waiting consumer must not starve indefinitely. |
| **Re-plan on change** | A material change in device state, context/KV, placement, backend or competing allocation invalidates the fit and requires re-planning. |
| **No silent substitution** | A failed or unavailable capability does not silently become another model, backend or placement. |
| **Cross-instance is explicit** | Two instances either share a declared authority or explicitly limit concurrent heavy work. A per-process ledger is not a global guarantee. |

### 7.3 Capability contract for future modalities

**HRA-33B (Proposed — required shape, not a final schema).** Every modality must declare enough for the governor to make an admission decision without hardcoding that modality into chat.

| Contract field | Purpose |
|---|---|
| Capability identity/version | Stable name and contract revision |
| Inputs/outputs/provenance | What it consumes and returns, with origin and confidence |
| Resource envelope | Weights, runtime overhead, context/KV, peak workspace, disk, CPU/GPU/thermal class |
| Placement support | GPU, CPU, mixed, multi-GPU, backend-specific limits |
| Lifecycle | Load, warm, probe, drain, unload, release evidence |
| Readiness probe | Workload-specific proof, not process-start proof |
| Concurrency | Allowed parallel requests, queue class, interruptibility |
| Priority/latency | Interactive/background, deadline sensitivity, fairness class |
| Workspace/storage | Temporary files, cache, cleanup, free-space requirement |
| Privacy/local-processing | Where data goes and what may be persisted |
| Failure modes | OOM, telemetry loss, thermal, driver loss, timeout, cancellation |
| User-work preservation | What must be snapshotted or retained across transition |

**HRA-33C (Calculated).** A capability that cannot declare its envelope is not ready for admission; it may run only under an explicit conservative policy that says so.

---

## 8. Evidence ladder for all state transitions

**HRA-34 (Owner-approved principle).** “Ready” is an evidence claim, not a status label. The same distinction applies beyond model loading:

- a chat turn is not persisted because `SaveChat` returned;
- a modality is not ready because its process started;
- a state transition is not complete because the UI changed;
- a memory write is not committed because the frontend called a save function;
- a model is not released because an unload call returned.

**HRA-35 (Proposed — general ladder).**

| Level | Meaning | Permitted claim |
|---|---|---|
| **Requested** | Intent was recorded | “The user asked for this” |
| **Admitted** | Policy and capacity allowed an attempt | “Howl decided to try” |
| **Attempted** | The operation was dispatched | “The operation was started” |
| **Observed** | A result or state change was seen | “Something appears to have happened” |
| **Verified** | A workload-appropriate check confirms the required outcome | “This outcome is established for this scenario” |
| **Unknown / failed** | Evidence is missing or contradictory | “We do not know” or “this failed” |

**HRA-36 (Calculated — no promotion by adjacency).** A process start does not promote readiness. A save call does not promote persistence. A UI update does not promote state transition. A log line does not promote release. Each claim needs its own evidence.

**HRA-37 (Proposed — evidence carries provenance).** Every consequential claim should retain source, timestamp, confidence and scope. “Measured,” “estimated” and “unknown” must remain distinguishable in the record and the UI.

---

## 9. Traceability chain

**HRA-38 (Proposed — required breadcrumb).**

The planning record should make this chain traversable:

**Original requirement → recovered evidence → owner decision → architecture constraint → design decision → work package → acceptance test → implementation evidence.**

| Link | Current artifact(s) | Status |
|---|---|---|
| Original requirement | `PROJECT_ISSUES_AND_EXPANSION.md`; `ISSUES.md`; original T1.2 and T0.5 | Source-derived; preserved |
| Recovered evidence | `HFS-AUDIT-2026-11-21-01`; `HFR-CP-2026-11-21-01` | Audits, not implementation |
| Owner decision | Current modality table; handoff requirement; “selection is not loading”; “ready is evidence” | Owner-approved, dated 2026-11-21 |
| Architecture constraint | This baseline, §3–§8 | Proposed/owner-approved as labeled |
| Design decision | `AMENDMENT-01`; `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`; unresolved `DQ-HRA-*` | Prior proposals plus open questions |
| Work package | HCDM M4 and future HRA gates; no implementation authorized | Proposed only |
| Acceptance test | This baseline §10; future HRA acceptance record | Not run |
| Implementation evidence | None for this baseline | Unestablished |

**HRA-39 (Calculated — no silent authority).** A later architecture document may refine this baseline, but it must cite which owner decision it implements and which unresolved question it closes. It must not silently reverse a scope decision or promote a proposal to an owner requirement.

---

## 10. Acceptance requirements for the future architecture gate

**HRA-40 (Proposed — not run).** These are future acceptance targets, not results.

1. **Selection ≠ loading:** selecting a model changes selection only; no allocation occurs until admission.
2. **Admission is atomic:** two concurrent requests cannot reserve the same free bytes.
3. **Coexist or transition:** safe coexistence is allowed when evidence supports it; otherwise the system queues, transitions, asks or refuses explicitly.
4. **Drain:** a busy capability is drained or refused, never silently killed.
5. **Release evidence:** unload reports verified, partial or unknown release; unknown is never shown as free.
6. **Readiness evidence:** a loaded process is not dispatched to until a workload-appropriate probe passes.
7. **Failure recovery:** failed load/readiness leaves an explicit state, preserves selection and user work, and does not silently substitute.
8. **Work preservation:** unsaved prompts, settings, chat state and pending work survive a handoff or are explicitly marked interrupted.
9. **Uncertainty policy:** missing telemetry produces a visible, conservative, defined behavior; it is not converted to certainty.
10. **Manual override:** soft limits can be overridden with consequences shown and recorded; hard safety limits remain non-overridable.
11. **Cross-instance:** two instances either coordinate through a declared authority or explicitly limit concurrent heavy work; they do not double-spend capacity silently.
12. **Orphan handling:** an orphan is reported and recoverable; it is not killed merely because it is unexplained.
13. **Scope boundary:** image generation, vision, video, audio input, audio generation and TTS use capability contracts; none is silently reduced to a text-model special case.
14. **Evidence ladder:** persistence, readiness, release and UI state transitions each have a verified claim distinct from requested/attempted/observed.

**HRA-41 (Calculated — no text-only implementation authorization).** A text-only slice is a possible later implementation strategy. It is **not** an approved design decision. The resource model, lifecycle protocol, ownership/concurrency rules, failure/recovery behavior, evidence requirements and future-modality contracts must be established first.

---

## 11. Rejected alternatives

**HRA-42 (Rejected).**

| Alternative | Why rejected |
|---|---|
| Selection immediately loads the model | Conflates intent with allocation; surprises the user and races other work |
| Unload request equals release | Runtime, driver, cache and process ownership can retain memory |
| Fixed sleep after unload | A sleep is not evidence |
| Unknown telemetry treated as free/zero | Produces unsafe admission and false confidence |
| One universal one-model limit | `AMENDMENT-01` requires supported coexistence; serialization is a policy, not a ban |
| Chat owns hardware policy | Couples conversation logic to runtime details and creates multiple authorities |
| Model manager owns global safety | Lifecycle mechanics cannot prove whole-machine safety |
| Path-based process reaping | Can kill another instance's child; T0.5's job-object direction is safer |
| Automatically kill a busy model to satisfy a new request | Destroys user work and hides contention |
| Silent model/placement substitution | Breaks selection, quality expectations and provenance |
| A reservation ledger is treated as physical truth | The OS/driver/other applications can invalidate it |

---

## 12. Planning gates

**HRA-43 (Proposed — planning only; no implementation authorization).**

| Gate | Question | Deliverable | Implementation allowed? |
|---|---|---|---|
| **HRA-0** | Is the discovery baseline complete and traceable? | This record | No |
| **HRA-1** | What is the resource model and reservation contract? | DQ-HRA decisions; logical/physical semantics; cross-instance choice | No |
| **HRA-2** | What is the lifecycle and handoff protocol? | State machine, drain/release/load/probe/recovery contracts | No |
| **HRA-3** | What is the evidence and telemetry contract? | Evidence ladder, confidence/provenance, release/readiness probes | No |
| **HRA-4** | How do multiple instances coordinate? | Chosen authority model and failure behavior | No |
| **HRA-5** | What acceptance proves the architecture? | Fault-injection and end-to-end acceptance design | No |
| **HRA-6** | Is a bounded implementation slice authorized? | Separate owner authorization | Not yet |

**HRA-44 (Calculated — sequencing).** This architecture is foundational to the eventual text, image, vision, video, audio and memory workloads. It does **not** preempt M0A release truth or Phase 1 identity/context work. Its design can proceed in parallel, but implementation requires its own gate.

---

## 13. Open questions

**HRA-45 (Unresolved).**

| ID | Question |
|---|---|
| `DQ-HRA-01` | Which component is the sole writer of reservation state? |
| `DQ-HRA-02` | Can operations share a parent reservation, and under what sub-budget rules? |
| `DQ-HRA-03` | What is the required response to an estimate overrun? |
| `DQ-HRA-04` | Which reservation records survive a crash, and how are they reconciled? |
| `DQ-HRA-05` | What evidence distinguishes a stale reservation from an orphaned allocation? |
| `DQ-HRA-06` | How does Howl respond when another application takes memory after admission? |
| `DQ-HRA-07` | Is multi-instance concurrent heavy work supported, and through which authority? |
| `DQ-HRA-08` | Which readiness probe is valid for each capability and workload? |
| `DQ-HRA-09` | Which release evidence is obtainable on Windows per runtime/driver? |
| `DQ-HRA-10` | What are the non-overridable safety limits, as distinct from warned soft limits? |
| `DQ-HRA-11` | What is the user-work snapshot/restore contract across a handoff? |
| `DQ-HRA-12` | How are future modality capability envelopes declared and versioned? |

---

## 14. Verification record

**HRA-46 (Observed).**

| Check | Result |
|---|---|
| Production code/schema/API/migration changes | **0** |
| Test changes/executions | **0** |
| Models/GPU/process/benchmark runs | **0** |
| Protected/canonical records changed | **0** |
| This file | New, untracked discovery baseline |
| Authorization | **None** |
| Canonical hashes | r33 `dfb28b60…` · v9 `698254ca…` · ledger `e5008020…` · map `b4ba33cc…` |

**HRA-47 (Observed — sources).**

| Topic | Path |
|---|---|
| Existing hardware authority | `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` |
| Owner hardware amendment | `docs/architecture/AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md` |
| Original scope/evidence | `docs/architecture/HOWL-FOUNDATION-ORIGINAL-SCOPE-AND-EVIDENCE-CHAIN-2026-11-21.md` |
| Readiness | `docs/architecture/HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md` |
| Program map | `docs/architecture/HOWL-CHAT-CROSS-WORKSTREAM-DEPENDENCY-AND-DELIVERY-MAP-2026-11-21.md` |
| Phase 1 | `docs/architecture/HOWL-FOUNDATION-PHASE-1-CHAT-STATE-CONTEXT-INTEGRITY-2026-11-21.md` |
| Process ownership | `internal/proc/job_windows.go`; `d3_chat_ownership_test.go`; `internal/proc/job_windows_test.go` |

**End of baseline.** The resource system is not a model manager with warnings. It is the application's authority for deciding what may run, proving what was released, proving what is ready, preserving user work, and recording the difference between a request, an attempt, an observation and a verified outcome.
