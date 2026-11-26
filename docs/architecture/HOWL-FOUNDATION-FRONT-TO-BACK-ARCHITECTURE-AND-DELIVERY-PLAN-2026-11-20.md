# Howl Chat — Front-to-Back Foundation Architecture and Delivery Plan

## Planning addendum `HF-ARCH-PLAN-2026-11-20-01`

**Status:** Planning and DeepSeek-paperwork preparation only. This addendum makes no implementation authorization, architecture adoption, schema/API decision, migration decision, test result, performance claim, or release claim.

**Change discipline:** This is a new, dated addendum. It does not alter the historical audits, owner records, frozen artifacts, decision registers, or the modified document registry. It must be read with the source records it links below.

**Scope:** One reconciled system plan for Howl Chat: conversation state, context construction, multi-character isolation, persistence and authoring, model and modality lifecycle, hardware admission, recovery, delivery gates, and evidence. It is deliberately front-to-back before work decomposition.

**Evidence labels:**

| Label | Meaning |
|---|---|
| **Verified in current tree** | Direct static inspection in this planning pass; not runtime proof. |
| **Prior reported evidence** | A result reported by a prior record, not rerun here. |
| **Owner-approved requirement** | A user/owner direction recorded in the supplied scope or a controlling owner record. |
| **Proposal** | A recommended design, not a decision. |
| **Open decision** | Requires explicit owner/architecture selection before implementation. |

## 1. Authority, traceability, and reconciliation

### 1.1 Controlling inputs

| Input | Role in this plan | Status retained |
|---|---|---|
| `PROJECT_ISSUES_AND_EXPANSION.md`, `ISSUES.md` | Recovered original product scope: the 27-item union, including T0.5. | Source-derived, not a signed baseline. |
| `docs/architecture/HOWL-FOUNDATION-ORIGINAL-SCOPE-AND-EVIDENCE-CHAIN-2026-11-21.md` (`HFS-AUDIT`) | Original-scope/evidence reconciliation and historical status correction. | Read-only audit. |
| `docs/architecture/HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md` (`HFR`) | Current readiness categories, test/result limits, and MEMORY-09 separation. | Readiness audit only. |
| `docs/architecture/HOWL-FOUNDATION-PHASE-1-CHAT-STATE-CONTEXT-INTEGRITY-2026-11-21.md` (`HF1`) | Bounded identity/context work-package and acceptance direction. | Proposed/authorization request. |
| `docs/architecture/HOWL-CHAT-CROSS-WORKSTREAM-DEPENDENCY-AND-DELIVERY-MAP-2026-11-21.md` (`HCDM`) | Cross-workstream order and M0A–M8 map. | Planning map. |
| `docs/architecture/HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md` (`M0A`) | Immediate release-truth/proc-safety acceptance review. | Criteria only; not executed. |
| `docs/architecture/HOWL-RESOURCE-ARBITRATION-AND-VERIFIED-HANDOFF-DISCOVERY-BASELINE-2026-11-21.md` (`HRA`) | Corrected modality ledger and resource/handoff discovery. | Discovery baseline and open decisions. |
| `docs/architecture/REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md` (`HSL`) | Ownership split: Howl owns chat lifecycle; MEMORY-09 owns memory consequences. | Owner-stated requirement; contract deferred. |
| `docs/architecture/AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md` and `AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` | Existing hardware and continuity requirements/design context. | Do not silently supersede. |

### 1.2 Reconciliation rules

1. The recovered original scope remains traceable as 27 active T0–T3 items. T0.5 remains open unless M0A proves it closed; its later omission is not closure.
2. A static code reading is not runtime, performance, UI, device, or acceptance evidence.
3. A prior DeepSeek report, test name, worklog, or historical live check is retained as **prior reported evidence** until rerun against an agreed gate.
4. The current modality ledger supersedes earlier scheduling language only where it explicitly says so: audio generation complete; TTS deferred but required for chat finalization; image generation in scope and nearly complete; vision largely complete but assessed for gaps; video in scope and incomplete; 2D/3D avatars out of scope.
5. Howl and MEMORY-09 remain separate programs. This plan specifies only the future boundary; it does not import MEMORY-09 rules into Howl or make MEMORY-09 implementation work.
6. Every delivery artifact must preserve: **source requirement → decision → contract → work package → acceptance evidence → observed result**.

## 2. Verified baseline of the present implementation

### 2.1 Conversation, context, and persistence

**Verified in current tree.** The Go/Wails application exposes chat bindings in `app-chat.go`. It keeps an in-memory mutable `a.chatMessages` list, sends through the llama HTTP client, and asynchronously calls `performSave` after a successful non-stream response. `recordUserTurn` contains a consecutive-duplicate guard. The streaming and frontend paths are distinct.

**Verified in current tree.** `internal/backend/chat/service.go` owns a second chat service/context. Its current send paths append the user message, use configured history/token trimming, resolve lore, and choose `characterNames[0]` as the active character. It is not yet a multi-character state authority.

**Verified in current tree.** `internal/backend/chat/persistence.go` has a per-chat directory with `meta.json`, a snapshot `history.json`, `history.ndjson`, world/scenario/character/lorebook/memory snapshots, and configuration pointing to the last chat. Load treats several data files as optional/best-effort. This is useful recovery material, but it is not a defined branch/lineage store.

**Verified in current tree.** Frontend code still has DOM-to-history synchronization and local-storage enrichment/recovery paths for characters and lorebooks. Therefore frontend state, backend session data, chat snapshots, append logs, content-library stores, and browser local storage are not yet governed as one authoritative composition path.

**Verified in current tree.** The present token budget is not end-to-end trustworthy: at least one prompt path estimates by characters-per-token and the backend has placeholder trimming. A real model/tokenizer measurement boundary remains required.

### 2.2 Multi-character and authoring ownership

**Verified in current tree.** Lorebook routing recognizes stable owner IDs and has privacy/conflict tests. World and scenario canonical loading uses `worldstore` in parts of the backend. Character loading and frontend authoring still include legacy/path/local-storage routes.

**Verified in current tree.** The frontend contains FIFO multi-character orchestration, but the inspected backend send path selects the first configured character. No direct static proof was found of per-character context isolation, state flushing, or a contamination test across speakers.

**Verified in current tree.** Chat persistence serializes composed copies of world, scenario, characters, lorebooks, and memories. That is not itself proof of canonical authoring ownership, revision pinning, or synchronization correctness.

### 2.3 Resource, model lifecycle, and process ownership

**Verified in current tree.** `internal/resource.Manager` tracks model registrations and can refresh a hardware snapshot. It is not a durable reservation/admission authority with cross-instance atomic claims.

**Verified in current tree.** `app-vram-arbitration.go` checks LLM/image estimated fit against refreshed VRAM and can unload one for the other, sleeps, then rechecks. It only covers a specific LLM/image handoff; it does not establish general release verification, lifecycle states, policy, other modalities, or cross-instance safety.

**Verified in current tree.** Shutdown is coordinated, cancels active streams/jobs, attempts final chat save, stops owned LLM/image/audio resources, and closes a Windows process job. This is a useful ownership base. It does not by itself prove two independently launched Howl instances cannot interfere or double-spend device capacity.

### 2.4 Modalities

| Capability | Current scope/status | Planning consequence |
|---|---|---|
| Text chat | Foundational and incomplete in state/context proof. | Build its authoritative state/context contract first. |
| Audio generation | Owner-reported complete. | Keep as a capability in the common lifecycle/resource contract; seek closure evidence separately. |
| TTS | Deferred, but required for chat finalization. | Do not hide it in a later studio; specify it now and gate final chat release on its approved acceptance. |
| ASR/audio input | Separate capability with substantial code/test surface. | It must have its own input/privacy/readiness contract. |
| Image generation | In scope; nearly complete. | Preserve its progress; move it under common admission/lifecycle without reopening finished UI work. |
| Vision | Largely complete; capability detection exists. | Audit remaining product integration and evidence gaps before closing. |
| Video | In scope and incomplete. | Treat as a future capability, not an image-engine afterthought. |
| 2D/3D avatars | Out of scope. | Exclude from architecture and acceptance gates. |

## 3. Target architecture: one system, explicit authorities

### 3.1 Boundary map

```text
UI / Wails bindings
  → Application command layer (validation, idempotency, user-visible outcomes)
    → Conversation authority (chat lineage, turns, current state, durable events)
      → Composition service (world/scenario/character/user/lore revisions → prompt plan)
        → Context-budget service (real tokenizer, allocations, renderable prompt)
          → Capability broker (admit / reserve / queue / handoff / recover)
            → Capability adapters (LLM, image, vision, video, ASR, TTS, audio generation)
              → owned runtime processes, OS/driver telemetry, model files

Authoring stores ───────────────→ composition service
Conversation event store ───────→ restore/replay and MEMORY-09 lifecycle notices
Evidence store ────────────────→ UI status, audit receipt, acceptance records
```

No lower layer may invent an upstream fact. UI rendering does not establish a committed turn; a model process does not establish readiness; a registry estimate does not establish free VRAM; a memory subsystem does not infer a chat transition.

### 3.2 Authorities and non-authorities

| Concern | Sole authority (proposal) | Explicit non-authorities |
|---|---|---|
| Current chat state, branches, deletion | Conversation authority | UI, prompt builder, MEMORY-09, append-log reader |
| Canonical authored objects and revisions | Content-library authority | chat snapshot, browser local storage, display names |
| Composition snapshot for a turn | Composition service | individual UI panels, raw object folders |
| Prompt/token allocations | Context-budget service with selected tokenizer | character count heuristic, renderer |
| Admission/reservation state | Capability broker | chat service, model adapter, per-modality manager |
| Process ownership/termination | Owning adapter and process/job registry | path-name reaper, another application instance |
| Readiness/release evidence | Capability-specific probe/evidence service | process-start event, sleep, UI label |
| Memory consequence of a declared chat event | MEMORY-09 under the deferred contract | Howl’s conversation authority |

### 3.3 Durable identities

The implementation paperwork must define immutable IDs, not choose them implicitly in code:

| Object | Required identity relationship |
|---|---|
| Content root | Identifies the authority/store namespace. |
| World, scenario, character, user persona, lore item | Stable object ID plus immutable revision ID/content hash. Display names are not keys. |
| Conversation | Stable conversation ID and root/store identity. |
| State/branch | Stable state ID, parent state ID(s), and explicit current-state selection. |
| Turn/message | Stable turn ID/message ID, author/character ID, source state, ordinal only as presentation metadata. |
| Prompt plan | Immutable plan ID bound to exact state, content revisions, tokenizer, model template, and budget. |
| Capability/model | Capability ID, implementation/version, model artifact identity, placement, process/job identity. |
| Reservation/job | Lease/job ID, owner instance ID, capability ID, snapshot/evidence IDs, timestamps. |
| Evidence receipt | Claim ID, operation ID, source, observation time, scope, confidence, and verifier version. |

## 4. Required state machines and contracts

### 4.1 Conversation state machine

**Owner-approved constraints from HSL:** Howl owns state selection/branching/regeneration/deletion; non-destructive transitions retain material; deletion is the only destructive operation; MEMORY-09 receives explicit facts and never infers transitions.

**Proposal:** Model conversation as an immutable event stream plus materialized current projection.

| Command | Preconditions | Durable result | Required notice/evidence |
|---|---|---|---|
| Create conversation | Valid root/content scope | Root state and initial projection | Creation receipt |
| Append user turn | Current state is writable; composition available | Turn event and successor state | Commit receipt; state-current notice |
| Generate response | Prompt plan accepted; model admitted | Attempt record; on completion assistant turn and successor state | Dispatch/readiness/commit receipts |
| Generate alternative | Base state retained | Sibling successor state; no overwrite | Alternative notice |
| Regenerate | Target response identified | New successor state; old material retained | Regeneration notice |
| Select/revert | Target retained state exists | Current-state pointer only | Selection notice |
| Edit | Target and edit semantics explicit | New derived state; prior material retained | Edit notice |
| Delete | Explicit confirmed deletion set | Tombstone/deletion event; no silent rewrite | Named deletion-set notice and report |
| Restore | Retained/tombstoned material permitted by policy | New current selection or recovery branch | Restoration receipt |

The materialized snapshot is a cache/recovery aid. It must be reconstructible from the event stream and must never silently replace it. `history.ndjson` can inform migration design, but no migration is authorized by this plan.

### 4.2 Turn and streaming state machine

`Draft → Submitted → Validated → Composed → Budgeted → AwaitingAdmission → Admitted → Generating → Completed | Failed | Cancelled → Persisted`.

Rules for paperwork:

1. User input is assigned its turn ID before inference; retries reuse an idempotency key and never append a second logical turn.
2. A streamed display delta is provisional until final completion and durable commit are separately reported.
3. A failed/cancelled generation preserves the user turn, prompt-plan identity, and failure evidence without manufacturing an assistant completion.
4. State-changing UI commands are serialized per conversation, but capability execution and other conversations may be concurrent when admitted.
5. “Persisted” has a defined durability level. The acceptance plan must state whether a write/flush/atomic-replace/reopen check is required for each critical record.

### 4.3 Composition and context contract

For each generation, the composition service emits a **PromptPlan**, not a concatenated anonymous string:

```text
PromptPlan {
  conversation/state/turn IDs;
  model-template + tokenizer identity;
  world/scenario/character/user-persona IDs and revisions;
  selected lore/memory references, owner/visibility/trigger decisions;
  ordered segments with byte/token counts and budgets;
  transcript selection and branch path;
  modality attachments and provenance;
  total budget, reserve-for-output, overflow decision;
  content hashes and policy version.
}
```

The context-budget service must count the rendered prompt with the selected model tokenizer, reserve output capacity, account for template and multimodal overhead, and return an explicit refusal/truncation/summarization policy outcome. It may not quietly discard private character material or infer costs from string length.

### 4.4 Multi-character isolation

For a multi-character scene, each speaker has a distinct **CharacterExecutionContext** bound to its character ID, private content scope, public/shared scope, selected state, prompt plan, token allocation, and output turn ID. The orchestrator schedules the next speaker; it does not build a shared flat prompt and relabel the result afterward.

Minimum isolation invariants:

- A private lore/memory record is ineligible unless its stable owner/visibility policy admits it to that speaker’s plan.
- A character receives a rendered transcript projection, never another character’s mutable working context/KV by accident.
- A speaker switch invalidates/flushes adapter context unless a verified, keyed reuse contract proves it is safe.
- Each response records which character/context/template/model generated it.
- Queue order is a scheduling choice, never identity or authority.

### 4.5 Authoring and synchronization contract

Worlds, scenarios, characters, user personas, and lorebooks need one canonical authoring API/store per object class. Chat binds immutable **references** (object ID + revision), not mutable copies as authority. It may snapshot the references for reproducibility.

The paperwork must name each existing bypass (frontend local storage, legacy folders, direct save routes, chat snapshots), classify it as read-only migration input, temporary compatibility view, or prohibited write path, and give it a retirement/containment gate. Display names cannot resolve ambiguous owners or replace IDs.

### 4.6 Capability lifecycle and resource admission

Every modality implements a common contract: declared envelope, compatible placements, load/warm/probe/drain/unload operations, concurrency class, input/output provenance, cancellation behavior, workspace cleanup, user-work preservation, and workload-specific readiness/release evidence.

`Selected → Planned → Reserved → Loading → Probing → Ready → Busy → Draining → Releasing → Released` with terminal/repair states `Refused`, `Unavailable`, `Failed`, `Unknown`, and `Recovering`.

Rules:

1. Selection does not allocate.
2. Planning uses model/profile estimates plus current telemetry and records uncertainty.
3. Reservation is logical; physical telemetry can veto it at commit.
4. A busy capability drains or refuses according to policy; it is not silently killed for a new request.
5. Load success is not ready. Unload success is not release. Both require their designated evidence.
6. Failure preserves selection and user work; it never silently substitutes a model, placement, or provider.
7. A crash makes leases untrusted. Recovery reconciles ledger, owned processes/jobs, and fresh telemetry before reclaiming.

### 4.7 Multi-instance model

This is an open architecture decision, not a test detail. Before concurrent heavy work is claimed safe, select one of: a single-supported-heavy-instance policy, shared local broker, atomic local lease registry, or a conservative explicitly non-guaranteed partition. Process job ownership solves child ownership, not global VRAM double-spending. M0A can prove non-interference in its narrow process case; it cannot close general resource arbitration.

## 5. Delivery architecture and gates

| Gate | Outcome | Depends on | Exit evidence |
|---|---|---|---|
| **M0A** Release truth/process safety | T0.1/T0.2/T0.5 truth in the current surface. | None; narrow first slice. | M0A A1–A9 including async append ordering and two independent instance tests. |
| **M0B** Ownership/ghost truth | Inventory of active/bypass/ghost UI and ownership paths. | M0A. | Zero selected ghost writes; canonical/bypass map and regression seams. |
| **M1** Conversation identity/restore | Selected state, IDs, durable event/snapshot recovery. | M0B; owner decisions on lineage/retention. | Append/save/restart/branch/revert/restart preservation scenarios. |
| **M2** Context composition/budget | Reproducible prompt plans and real token accounting. | M1; tokenizer/template policy. | Exact plan receipt; known over-budget/refusal/truncation cases; no hidden-thought leakage. |
| **M3** Multi-character isolation | Per-speaker contexts and safe switches. | M1/M2; authoring ownership map. | Conflicting-private-fact scene, queue order, no cross-character injection, flush/reuse evidence. |
| **M4** Resource/lifecycle authority | Common capability contract, admission, handoff, recovery. | HRA decision set; M0A process ownership. | Atomic reservation, verified release/readiness, failure restore, cross-instance policy acceptance. |
| **M5** Authoring/lore/UI completion | Canonical authoring and remaining original product work. | M0B–M3. | Object/revision synchronization; remaining T2/T3 acceptance by item. |
| **M6** Performance/compatibility | Measured, reproducible performance and hardware tiers. | M2/M4. | Pre-registered benchmark matrix; no estimate presented as measurement. |
| **M7** ASR end-to-end | Audio input as accepted chat capability. | M2/M4 and explicit privacy/consent behavior. | UI-to-transcript-to-turn test with device and recovery cases. |
| **M8** MEMORY-09 readiness | Interface readiness, not product-memory integration. | M1/M2 and HSL contract decision. | Contract fixture/replay evidence; Howl/MEMORY-09 responsibilities remain separate. |
| **Chat finalization** | TTS plus accepted chat foundation. | M1–M4 and TTS-specific gate. | Text-to-voice handoff, cancellation, playback, persistence/privacy, lifecycle evidence. |
| **Image/Vision/Video finalization** | Modalities meet their own product gates. | M4 and applicable composition/UI gates. | Modality-specific acceptance; no avatar work. |

Parallelism is allowed only where contracts do not change concurrently: M0A/M0B evidence work can proceed while HRA paperwork is completed; M4 design can proceed alongside M1–M3 design; modality-specific adapter audits can proceed under stable capability-contract drafts. Implementations that consume changing identity, resource, or authoring contracts must wait for their gate.

## 6. Workstreams for DeepSeek paperwork

1. **Conversation foundation:** data model alternatives, transition table, event/snapshot durability, command idempotency, restoration and deletion semantics. Inputs: HSL, HF1, M0A.
2. **Context/composition:** canonical-reference resolution, ownership/visibility, tokenizer/template interface, segment budgets, reproducible prompt-plan receipt.
3. **Multi-character:** speaker isolation, transcript projection, context/KV flush or keyed reuse proof, queue policy, attribution.
4. **Authoring and content ownership:** authoritative stores and revisions; frontend/backend synchronization; inventory and disposition of every bypass/legacy route.
5. **Resource/capability platform:** broker, reservation model, telemetry normalization, lifecycle adapter protocol, cross-instance choice, recovery and evidence receipts.
6. **Modalities:** separate capability papers for image, vision, video, ASR, TTS, and audio generation, all conforming to the shared contract.
7. **Reliability and performance:** fault matrix, recovery protocol, test fixture ownership, benchmarks, compatibility tier evidence.
8. **MEMORY-09 boundary:** event-notification schema and replay/idempotency semantics only after ownership decisions; no memory-rule redesign.

Each paper must include: purpose/scope, source IDs, invariants, authority/non-authority table, state transitions, request/response/error contracts, persistence/durability, telemetry/evidence, recovery, security/privacy, backward-compatibility/migration questions, acceptance tests, observability, open decisions, and explicit implementation authorization boundary.

## 7. Acceptance and evidence strategy

### 7.1 Evidence ladder

Every consequential claim is labelled **requested → admitted → attempted → observed → verified**, or **unknown/failed**. UI state, process start, log output, and function return cannot advance a claim without the evidence designated for that transition.

### 7.2 Test layers

| Layer | Purpose | Examples |
|---|---|---|
| Unit/property | Pure state and contract laws. | IDs stable; branch transitions retain material; reservation cannot be double-spent. |
| Adapter/component | Runtime boundary behavior with controlled fakes. | Load/probe/drain/release, telemetry loss, process ownership. |
| Persistence/recovery | Crash-safe reconstruction and idempotent replay. | Mid-stream exit; incomplete commit; stale lease; corrupt snapshot plus valid event log. |
| Integration | Cross-service truth. | Prompt plan uses correct revisions and tokenizer; model handoff preserves unsent work. |
| End-to-end | User-observable acceptance. | M0A A1–A9, multi-character privacy scene, TTS turn playback/cancel, image/video workflows. |
| Performance | Measured capacity/latency behavior. | Pre-registered hardware/model/modality matrix; cold/warm, contention, long-context cases. |

No stress, model, device, or destructive run is authorized by this planning artifact. Each later test package must declare cleanup, data roots, process ownership, expected writes, and stop conditions before execution.

### 7.3 Required receipts

- Conversation command/transition receipt with IDs and durability result.
- Prompt-plan receipt with source revisions, segment counts, tokenizer/template, and budget decision.
- Admission/reservation receipt with snapshot, estimate, policy, uncertainty, and lease owner.
- Capability lifecycle receipt with load/readiness/release evidence.
- Recovery receipt with discovered durable state, process/job reconciliation, and disposition.
- Acceptance receipt with build identity, fixtures, exact scenario, artifacts, pass/fail/unknown, and limitations.

## 8. Risks and stop conditions

| Risk | Stop condition / mitigation |
|---|---|
| Frontend/backend/persistence disagreement | Stop any claim of correct restore until a single authoritative write/read path and reconciliation test exist. |
| Branch redesign leaks into M0A | Stop M0A; keep it limited to append/render/process ownership unless the owner authorizes M1. |
| Cross-character information leak | Stop multi-character feature acceptance; retain fixture and prompt-plan evidence; do not mask by UI filtering. |
| Hardware claim based on estimates or sleeps | Report unknown/not-ready; do not dispatch or assert safe release. |
| Two instances race for capacity | Enforce selected multi-instance policy before concurrent heavy work. |
| Legacy bypass mutates canonical data | Contain/disable the write path under an explicit compatibility decision; do not silently merge. |
| MEMORY-09 scope creep | Stop at lifecycle fact interface; do not alter its records, rules, or experiment. |
| Modalities treated as special cases | Reject integration until they supply the common capability contract and acceptance evidence. |

## 9. Decisions required before implementation

### 9.1 Conversation and authoring

| ID | Decision question | Blocking effect |
|---|---|---|
| DQ-C01 | Is the state unit a message alternative, branch, lineage, or another explicitly defined graph? | M1 storage/API. |
| DQ-C02 | What retained material is queryable, by whom, and for how long? | Branch/recovery/privacy. |
| DQ-C03 | What deletion/retention promise and confirmation semantics apply? | Destructive lifecycle and HSL boundary. |
| DQ-C04 | What durability level is required per transition before the UI reports committed? | Crash/recovery acceptance. |
| DQ-C05 | Which authoring store is canonical for each object, and what is each bypass’s disposition? | M0B/M1/M5. |
| DQ-C06 | How are legacy IDs, display names, and copied snapshots mapped/quarantined without inventing ownership? | Migration/recovery. |
| DQ-C07 | What are user-persona ownership, attachment, and active-identity rules? | Composition/UI contract. |

### 9.2 Context and multi-character

| ID | Decision question | Blocking effect |
|---|---|---|
| DQ-X01 | Which tokenizer/template source is authoritative for each model family? | Real budget proof. |
| DQ-X02 | What is the segment-budget policy under overflow: refuse, trim, summarize, ask, or model-specific alternative? | Prompt determinism. |
| DQ-X03 | Is context/KV always flushed on speaker switch, or can keyed reuse be admitted with evidence? | Isolation and latency. |
| DQ-X04 | What public/shared transcript projection is visible to each character? | Privacy/scene semantics. |
| DQ-X05 | What multi-character scheduling/fairness/cancellation policy applies? | Orchestration UX. |

### 9.3 Resources and modalities

| ID | Decision question | Blocking effect |
|---|---|---|
| DQ-R01 | Which component is the sole reservation writer? | M4 correctness. |
| DQ-R02 | Are multiple instances supported for concurrent heavy work, and which coordination authority is selected? | Cross-instance safety. |
| DQ-R03 | What is a hard limit versus an overrideable soft limit? | Admission UX/safety. |
| DQ-R04 | Which telemetry/release probes are trustworthy per Windows runtime/driver? | Verified handoff. |
| DQ-R05 | What survives crash and how are stale leases/orphaned processes reconciled? | Recovery. |
| DQ-R06 | What exact readiness probe is valid for each capability/workload? | Dispatch truth. |
| DQ-R07 | How do capability resource envelopes version and declare uncertainty? | Future modality integration. |
| DQ-R08 | What user work must be snapshotted/restored across handoff and failure? | No-loss guarantee. |

### 9.4 MEMORY-09 boundary

| ID | Decision question | Blocking effect |
|---|---|---|
| DQ-M01 | What ordered, idempotent lifecycle event contract carries state identity, affected material, and deletion set? | M8 interface. |
| DQ-M02 | What does Howl guarantee for replay, retention, and a MEMORY-09 unavailability/failure response? | Recovery and non-destructive semantics. |
| DQ-M03 | Which side owns applicability representation? | Joint contract only; no unilateral selection. |

## 10. DeepSeek paperwork sequence

1. Confirm this addendum’s labels, scope corrections, and open-decision list; record any owner corrections as a dated addendum.
2. Produce a **Foundation Contract Pack** for conversation identity/state, composition/token budgeting, content ownership, and resource/capability lifecycle. Do not write implementation plans before unresolved decision rows are answered.
3. Produce an **Acceptance and Evidence Pack** with fixtures, receipts, fault matrix, M0A execution design, and no-claim rules.
4. Produce **modality annexes** against the shared capability contract, beginning with TTS because it gates chat finalization while retaining image/vision/video/ASR/audio-generation scope.
5. Select the smallest authorized implementation slice only after the relevant contracts and acceptance gate are approved. The current candidate remains M0A, narrowly scoped as its review defines.

## 11. Current conclusion

The system should be developed as a composition of authoritative state, reproducible prompt planning, isolated character execution, and capability-based resource admission—not as incremental patches to the UI, a single mutable chat history, or isolated model managers. M0A is the immediate proof-of-release-truth gate. M1–M4 establish the foundations that make chat, TTS, image, vision, video, and audio workflows trustworthy together. No existing record is closed by this plan; all closure requires the specifically defined evidence at its delivery gate.

## 12. Planning-pass verification record

| Check | Result |
|---|---|
| Repository access | Verified at `C:\\Users\\StoneWolfPC\\HOWL_Chat`. |
| Current code inspected | Chat bindings/persistence/service, frontend chat paths, lore routing, resource manager, VRAM handoff, lifecycle/process ownership, modality/test inventories. |
| Documentation inspected | HFS, HFR, HF1, HCDM, M0A, HRA, HSL and architecture registry/context. |
| Code/schema/API/migration change | None. |
| Tests/builds/models/GPU/stress runs | None. |
| Existing documentation modified | None. |
| New artifact | This planning addendum only. |

