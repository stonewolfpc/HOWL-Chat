# HOWL + MEMORY-09 — Program Direction, Foundation Readiness, and Checkpoint

**Status:** program direction, recovery index and checkpoint. **Not** architecture approval, implementation
authorization, a completion claim, or a MEMORY-09 rule change.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Owner direction | Finish MEMORY-09 and raise Howl's foundation in parallel; verify the combined system under real creative workloads |
| Checkpoint | **`CP-2026-11-21-01`** |
| Recovery token | **`HPD-CP-2026-11-21-01`** — search the repository for this token |
| This file | `docs/architecture/HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md` |
| Current phase | Direction set; consolidated readiness audit **complete** (`HFR-CP-2026-11-21-01`); evidence work not started |
| Canon | r33, v9, owner ledger and conformance map unchanged |
| Map | **118** requirements |
| Authorization | No production API, schema, code, model fetch, corpus, G2 or BGE work |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 0. Recovery card — read this first

**If we lose our place, start here.**

1. Read **§2 — Recovery index** for the records that define the program.
2. Read **§8 — Checkpoint `CP-2026-11-21-01`** for the current state.
3. Read **§9 — Immediate next move** and the readiness audit's §1, §10 and §13 before starting any new work.
4. Do not infer current status from an older record. Records are dated; the checkpoint says what was true on
   **2026-11-21**.

**One-sentence state:** MEMORY-09's architecture is strong but not executed; Howl has working pieces but no
end-to-end readiness evidence; the readiness audit is complete, so the next task is evidence work on state identity,
context budget, resource behavior and a real-server smoke matrix — not more implementation by default.

**One-sentence destination:** one integrated local creative environment where autonomous memory, dynamic world/scenario/
character state, hardware-aware inference and persistent knowledge work together under real use, without cloud
dependency for ordinary operation.

---

## 1. Destination — what we are building toward

**HPD-01 (Observed — owner direction).** The destination is **not** an experiment optimized around the limitations of
the current application. It is a finished local creative environment, tested under the conditions it is meant to run
in.

| Capability | Intended result |
|---|---|
| Multi-character environment | Multiple characters coexist, speak distinctly, share scenes without knowledge or voice contamination |
| Autonomous memory | Capture, classification, lifecycle and retrieval work without author-supplied classes or manual review |
| Dynamic state | World, scenario and character state can change independently and be restored |
| Hardware-aware inference | VRAM/RAM pressure, model loading/unloading, context sizing and recovery are handled deliberately |
| Persistent knowledge | Memory survives context flushing, model swaps, restarts and long conversations |
| Local-first | Ordinary operation does not depend on a cloud service |
| Owner-led verification | Real creative workloads, not only synthetic fixtures, are the final acceptance environment |

**HPD-02 (Calculated — the destination is a claim about direction, not current completion).** No row above is asserted
as operational. §5 records the current readiness baseline; §8 records what is actually established.

---

## 2. Recovery index — the links that define this program

**HPD-03 (Observed — entry point).** This file is the program entry point. The links below are relative to
`docs/architecture/`; they are the records to read before changing direction.

### 2.1 Direction and current state

| Purpose | Record |
|---|---|
| This program direction and checkpoint | **this file** |
| MEMORY-09 current-state audit (2026-11-18; execution and coverage caveats) | [`MEMORY-09-CURRENT-STATE-AND-COMPLETION-AUDIT-2026-11-18.md`](../research/MEMORY-09-CURRENT-STATE-AND-COMPLETION-AUDIT-2026-11-18.md) |
| MEMORY-09 autonomous-classification plan (capture → classification → unresolved lifecycle) | [`MEMORY-09-AUTONOMOUS-CLASSIFICATION-EXPERIMENT-PLAN-2026-11-21.md`](../research/MEMORY-09-AUTONOMOUS-CLASSIFICATION-EXPERIMENT-PLAN-2026-11-21.md) |
| Production capture boundary trace and proposal | [`MEMORY-09-PRODUCTION-CAPTURE-BOUNDARY-PROPOSAL-2026-11-21.md`](../research/MEMORY-09-PRODUCTION-CAPTURE-BOUNDARY-PROPOSAL-2026-11-21.md) |
| G1 fixture contract | [`MEMORY-09-G1-CAPTURE-FIXTURE-CONTRACT-2026-11-21.md`](../research/MEMORY-09-G1-CAPTURE-FIXTURE-CONTRACT-2026-11-21.md) |
| G1 fixture results | [`MEMORY-09-G1-CAPTURE-FIXTURE-RESULTS-2026-11-21.md`](../research/MEMORY-09-G1-CAPTURE-FIXTURE-RESULTS-2026-11-21.md) |
| **Consolidated foundation-readiness audit** | [`HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md`](HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md) |
| **Phase 1 work package — chat state/context** | [`HOWL-FOUNDATION-PHASE-1-CHAT-STATE-CONTEXT-INTEGRITY-2026-11-21.md`](HOWL-FOUNDATION-PHASE-1-CHAT-STATE-CONTEXT-INTEGRITY-2026-11-21.md) |

### 2.2 MEMORY-09 canon and decisions

| Purpose | Record |
|---|---|
| Architecture v9 — ratified design baseline | [`MEMORY-09-ARCHITECTURE-v9-2026-11-17.md`](../research/MEMORY-09-ARCHITECTURE-v9-2026-11-17.md) |
| r33 stage contract — memory lifecycle and authority source | [`MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md`](../research/MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md) |
| 118-requirement conformance map | [`MEMORY-09-CONFORMANCE-MAP-2026-11-17.md`](../research/MEMORY-09-CONFORMANCE-MAP-2026-11-17.md) |
| Owner decisions / T1 ledger | [`MEMORY-09-OWNER-DECISIONS-2026-11-17.md`](../research/MEMORY-09-OWNER-DECISIONS-2026-11-17.md) |
| L-2 — one formation entry, independent dimension precedence | [`MEMORY-09-L2-OWNER-DECISION-2026-11-20.md`](../research/MEMORY-09-L2-OWNER-DECISION-2026-11-20.md) |
| Decision A — reconstructibility from the authoritative stream | [`MEMORY-09-DECISION-A-OWNER-DETERMINATION-2026-11-20.md`](../research/MEMORY-09-DECISION-A-OWNER-DETERMINATION-2026-11-20.md) |
| Formation/classification relationship — options and hybrid direction | [`MEMORY-09-FORMATION-CLASSIFICATION-RELATIONSHIP-2026-11-21.md`](../research/MEMORY-09-FORMATION-CLASSIFICATION-RELATIONSHIP-2026-11-21.md) |
| Implementation plan — files, schema, tests, stop conditions | [`MEMORY-09-IMPLEMENTATION-PLAN-2026-11-20.md`](../research/MEMORY-09-IMPLEMENTATION-PLAN-2026-11-20.md) |
| A2 ordering and event-schema reconciliation | [`MEMORY-09-A2-ORDERING-AND-EVENT-SCHEMA-RECONCILIATION-2026-11-19.md`](../research/MEMORY-09-A2-ORDERING-AND-EVENT-SCHEMA-RECONCILIATION-2026-11-19.md) |

### 2.3 Howl foundation

| Purpose | Record |
|---|---|
| Chat-state lifecycle requirements (HSL-1…HSL-11, contract deferred) | [`REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md`](REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md) |
| Reconciled architecture and controlled design plan | [`RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`](RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md) |
| Integration and acceptance plan | [`INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md`](INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md) |
| Memory architecture decision | [`ARCHITECTURE-DECISION-MEMORY-SYSTEM-2026-11-12.md`](ARCHITECTURE-DECISION-MEMORY-SYSTEM-2026-11-12.md) |
| Memory continuity design | [`DESIGN-MEMORY-CONTINUITY-2026-11-11.md`](DESIGN-MEMORY-CONTINUITY-2026-11-11.md) |
| Hardware resource orchestration design | [`DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`](DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md) |
| Multi-character continuity and routing amendment | [`AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md`](AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md) |
| Foundational backend rebuild — decision record: five roots, configuration authority and routing record, persistence lifecycle, runtime resolution and evidence, ownership classes (D1–D5) | [`ADR-001-FOUNDATIONAL-BACKEND-REBUILD-2026-11-22.md`](ADR-001-FOUNDATIONAL-BACKEND-REBUILD-2026-11-22.md) |
| Storage locations, picker navigation and content routing — three state families, N1–N19; refines ADR-001 D1 for chats | [`ADR-002-STORAGE-AND-NAVIGATION-POLICY-2026-11-22.md`](ADR-002-STORAGE-AND-NAVIGATION-POLICY-2026-11-22.md) |

### 2.4 Code and test evidence

| Purpose | Path |
|---|---|
| G1 fixture harness (test-only) | [`internal/backend/memory/capture_contract_test.go`](../../internal/backend/memory/capture_contract_test.go) |
| MEMORY-09 authority engine (implementation, not connected to harness) | [`eval/mem1/m09/`](../../eval/mem1/m09/) |
| Conversation harness | [`eval/mem1/conversation.py`](../../eval/mem1/conversation.py) |
| Resource snapshot / allocation bindings | [`app-resource.go`](../../app-resource.go) |
| VRAM arbitration | [`app-vram-arbitration.go`](../../app-vram-arbitration.go) |
| Model intelligence | [`app-model-intelligence.go`](../../app-model-intelligence.go) |

**HPD-04 (Calculated — if a link is stale).** A stale record is not a reason to infer a new state. Re-read the source
and add a dated correction instead. The checkpoint is a snapshot, not a permanent truth.

---

## 3. Program shape — two workstreams, one integration seam

**HPD-05 (Calculated — parallel, not serial).**

```text
Workstream M — MEMORY-09
  capture → classification → unresolved lifecycle
  acquisition/lifecycle events → ordering/authority/provenance
  conformance → implementation/verification
                    │
                    ├── Howl ↔ MEMORY-09 contract
                    │
Workstream H — Howl foundation
  chat/state → multi-character → hardware
  inference → performance
                    │
                    ▼
        Integration verification under real workloads
```

**HPD-06 (Observed — no global barrier).** MEMORY-09 research and architecture can proceed while Howl's chat, state
and hardware work proceeds. Integration waits for the **specific interfaces**, not for every subsystem to finish.

**HPD-07 (Calculated — the production-boundary proposal is input, not instruction).** `PCB-01…42` describes how the
current app could deliver a capture envelope. It does **not** authorize changing the current application, and it does
not require the current application to remain the final shape.

---

## 4. Workstream M — finish MEMORY-09

**HPD-08 (Source-derived — current state).** MEMORY-09 is design-complete in its ratified rules and requirement map,
partially implemented, and largely unexecuted. The 2026-11-18 audit found no connection between the `m09` engine and
the conversation harness; the 2026-11-21 G1 fixture pass adds contract evidence, not production evidence.

### 4.1 Tracks

| Track | Goal | Current state | Gate / dependency |
|---|---|---|---|
| **M1 — Capture boundary** | Capture what happened, who experienced it, where it came from, with stable identity | G1 fixture contract **passes**; production boundary **proposed**, not implemented | Stable message/turn identity; separate production authorization |
| **M2 — Autonomous classification** | Classify only when capture is sufficient; otherwise preserve uncertainty | G2 **blocked**; no class-labelled corpus; rules/BGE design not executed | M1; corpus; G2 authorization |
| **M3 — Unresolved-class lifecycle** | Keep an unresolved memory useful without class-sensitive use | Analysis exists; **no L-2/Decision A amendment** | Owner decision; M1; no class invention |
| **M4 — Acquisition and lifecycle event semantics** | Make learning, telling, believing, forgetting and status changes explicit | A2 ordering selected; event-schema reconciliation exists; representation/implementation gaps remain | M5; integration contract |
| **M5 — Ordering, authority, provenance, reconstruction, integrity** | Complete the structural guarantees across streams | 57 rules and 118 requirements exist; **0 executed evidence**; some audit caveats remain | Conformance execution; reconciliation |
| **M6 — Conformance coverage** | Map every requirement to current evidence and expose gaps | 118 requirements defined; mapping not re-performed for the committed 1,951-line tree; **0 runtime** | M7; fresh mapping |
| **M7 — Implementation and verification** | Connect the engine, harness and production boundary; verify | `m09` exists but is not connected; tests not authorized beyond the fixture contract | Separate test authorization |
| **M8 — Integration contract** | Define Howl ↔ MEMORY-09 state, transition, deletion and failure semantics | HSL §7 deferred; no format, identifier scheme or protocol | HSL; M5; H foundation |

### 4.2 Remaining MEMORY-09 decision register

**HPD-09 (Calculated — initial register, not claimed exhaustive).** The readiness audit in §9 reconciles this with the
owner ledger.

| ID | Decision | Current position | Status | Blocks |
|---|---|---|---|---|
| **PD-M01** | Implement the production capture boundary | Server-side capture at append is recommended; no production change authorized | **Open** | M2, real-workload memory |
| **PD-M02** | Autonomous classification method | Deterministic rules first; hybrid/BGE candidate only; G2 blocked | **Open** | Class assignment |
| **PD-M03** | Unresolved-class lifecycle | Preserve uncertainty; no class-sensitive use; L-2/Decision A amendment required | **Blocked** | Retrieval/use of unresolved memories |
| **PD-M04** | Acquisition/lifecycle event representation | A2 ordering selected; exact event/status carrier not fully settled | **Partially decided** | M4, M5, integration |
| **PD-M05** | Ordering/authority/provenance/reconstruction/integrity closure | Design rules exist; execution evidence absent | **Open** | Acceptance |
| **PD-M06** | Conformance coverage and execution | 118 defined; 0 executed; mapping stale for current tree | **Open** | Implementation claims |
| **PD-M07** | Engine ↔ harness integration | No caller from harness to `m09`; no production integration | **Open** | Any meaningful memory test |
| **PD-M08** | Howl ↔ MEMORY-09 contract | HSL requirements recorded; contract deliberately deferred | **Open** | Integration verification |
| **PD-M09** | BGE embedding experiment | Separate model/corpus gate; no fetch; reranker is not an embedding model | **Deferred** | Classification candidate ranking only |
| **PD-M10** | Real-workload acceptance policy | Not defined; depends on H and M | **Deferred** | Final acceptance |

### 4.3 Non-negotiable MEMORY-09 boundaries

**HPD-10 (Observed — preserved).**

- r33 `scope` is **not** v9 authority `domain`.
- S2H-12: context inactivity is not a memory-lifecycle event.
- P-1/S-3 remain closed; Decision B is untouched.
- L-2 remains one formation entry with `scope_baseline`, `class_baseline`, `status_slot`; no universal or
  cross-dimension ordering.
- Decision A remains: initial `scope`, `knowledge_class`, `status` must be reconstructible from the authoritative
  stream; record fields are caches.
- Capture-incomplete is **not** classification-unresolved.
- No classifier may compensate for incomplete provenance.
- No character speech becomes authored-world knowledge through subject matter alone.
- No inferred `told` from a character merely speaking.

---

## 5. Workstream H — Howl foundation readiness baseline

**HPD-11 (Observed — this is a record-based baseline, not a finished code audit).** The states below are sufficient to
sequence work. §9 is the audit that converts them into pass/fail evidence.

### 5.1 Readiness matrix

| Area | What exists | What is missing / unverified | Baseline | Can proceed independently? |
|---|---|---|---|---|
| **A. Chat and state machine** | Chat history, message append, persistence, streaming, undo/edit paths; HSL-1…HSL-11 requirement record | Stable state identity; branch/alternative sets; revert; explicit lifecycle notifications; non-destructive retention; HSL §7 contract | **Not ready for memory integration** | **Yes** — this is Howl-owned |
| **B. Multi-character orchestration** | Character-specific response path with stable `CharacterID`; generic assistant path uses `"AI"` with no character ID | Response sequencing; context flushing; voice separation; shared-scene awareness; memory access boundaries; contamination tests | **Partial / unverified** | **Yes** — with A |
| **C. Hardware-aware intelligence** | Resource manager bindings; hardware snapshots; allocation advice; VRAM arbitration; model intelligence | One authoritative lease path; telemetry source/confidence; measured reserve; OOM prevention/recovery acceptance; provider order | **Partial / unverified** | **Yes** — with D |
| **D. Inference engine and Kobold-like capabilities** | Model registry, model discovery, streaming chat path, model load/unload and VRAM checks | Capability matrix for KV cache, context sizing, sampling, stop sequences, batching, speculative decoding, quantization; backend-specific evidence | **Partial / unverified** | **Yes** — with C |
| **E. Performance and hardware tricks** | Some measured constants and device-specific notes exist in code/design records | Reproducible benchmarks for layer placement, cache reuse, memory mapping, flash attention, async loading, parallelism; no assumption that an optimization helps | **Unestablished** | **Yes** — measurement only |

**HPD-12 (Calculated — readiness labels).**

| Label | Meaning |
|---|---|
| **Ready** | Implemented and verified for the named scope |
| **Partial / unverified** | Real code exists, but the required behaviour has no consolidated acceptance evidence |
| **Design-only** | Design or requirement record exists; no implementation claim |
| **Not ready** | A required capability is absent or contradicted by current behaviour |
| **Unestablished** | No sufficient evidence either way |

### 5.2 What the HSL record already establishes

**HPD-13 (Source-derived — `REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md`).**

| Requirement | Current gap |
|---|---|
| HSL-1: Howl owns the state machine | No state model exists |
| HSL-2: authoritative lifecycle notification | No interface exists |
| HSL-3: stable state identity | Message IDs are positional |
| HSL-4: alternative selection | No alternative set; Prev/Next is scrolling |
| HSL-5: regeneration | Re-streams without syncing history |
| HSL-6: revert | No revert action |
| HSL-7: deletion is the only destructive operation | Undo/delete rewrites history; previous state unrecoverable |
| HSL-8: deletion names its set | No deletion report |
| HSL-9: dependency determined by MEMORY-09 | Contract deferred |
| HSL-10: retention across non-destructive transitions | Not implemented |
| HSL-11: no unilateral decisions | No boundary exists yet |

**HPD-14 (Calculated — the HSL record is a requirement, not an implementation authorization).** It deliberately
authorizes no change. Workstream H must produce its own design and verification record before code changes.

---

## 6. Integration verification — the scenarios that matter

**HPD-15 (Observed — owner-defined scenarios).** These are the acceptance scenarios for the combined system, not
fixtures for one subsystem.

| Scenario | What must be verified | Prerequisites |
|---|---|---|
| **Character switches** | Correct memory activation; context isolation; no bleed between characters | M1, M8, H-A, H-B |
| **World/scenario switches** | Correct memory selection; eligibility and restoration; no cross-world leakage | M1, M5, M8, H-A |
| **Character forgets something** | Lifecycle history and projections remain correct; no fabricated retention | M3, M4, M5 |
| **Character learns something** | Acquisition, attribution and provenance are explicit; no inferred source | M1, M4 |
| **Long conversation** | Memory remains coherent as context is flushed or compacted | H-A, H-D, M1, M5 |
| **Model unload/reload** | Memory survives independently of inference state | H-C, H-D, M7 |
| **Multiple characters** | No knowledge or voice contamination; access boundaries hold | H-B, M1, M8 |
| **Hardware pressure** | Memory operations do not destabilize inference; no OOM spiral | H-C, H-E |
| **Rebuild/restart** | Authoritative history reconstructs the same state | H-A, M5, M7, M8 |

**HPD-16 (Calculated — no scenario is passed by a fixture alone).** The G1 fixture contract proves the data model can
express the envelope. It does not prove the app supplies it, the engine consumes it, or the combined system behaves
correctly under real use.

---

## 7. Sequencing and gates

**HPD-17 (Calculated — the program is dependency-driven, not date-driven).**

| Gate | Meaning | State |
|---|---|---|
| **G0** | Research protocol accepted | **Passed** |
| **G1** | Capture fixture contract | **Passed (fixture-only)** |
| **PCB** | Production capture boundary proposal | **Complete; not implemented** |
| **HFR** | Consolidated Howl foundation-readiness audit | **Complete (audit)** — see [`HFR`](HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md) |
| **M2** | Classification capability screen | **Blocked** on M1 and corpus |
| **H-A…H-E** | Foundation acceptance by area | **Not started** |
| **I** | Howl ↔ MEMORY-09 integration contract | **Deferred** |
| **RW** | Real-workload verification | **Not started** |

**HPD-18 (Observed — separate gates).**

- G2 here means the **G0–G5 research gate**, not Gate 2 implementation authorization.
- BGE remains separately gated: no model fetch, no corpus, no experiment until authorized.
- Production API/schema/code changes require separate authorization.
- A passing fixture does not authorize implementation.
- A design package does not authorize architecture approval.

**HPD-19 (Calculated — parallel work is allowed, integration claims are not).** Workstream H can make progress without
M; M can make research progress without H. Neither may claim the combined system works until the integration scenarios
in §6 have evidence.

---

## 8. Checkpoint `CP-2026-11-21-01`

**HPD-20 (Observed — resume here).**

| Field | State at this checkpoint |
|---|---|
| Program phase | Direction set; foundation-readiness audit complete; evidence work not started |
| MEMORY-09 architecture | Design baseline ratified; 57 rules; 118 requirements; r33/v9/ledger/map unchanged |
| MEMORY-09 implementation | `eval/mem1/m09/` exists (1,951 lines at the 2026-11-18 audit); not connected to the harness |
| MEMORY-09 execution | G1 fixture contract passes; **0 production/runtime memory verification** |
| Production capture | Trace and proposal complete; **no implementation** |
| Classification | G2 blocked; no class-labelled corpus; BGE gated |
| Unresolved-class lifecycle | Analysis only; L-2/Decision A amendment not authorized |
| Howl chat/state | HSL requirements recorded; no state machine/branch/revert implementation |
| Howl multi-character | Character path exists; orchestration and contamination acceptance unverified |
| Howl hardware | Resource manager/VRAM code exists; consolidated acceptance unverified |
| Howl inference/performance | Partial/unestablished; no target capability matrix or benchmark set |
| Integration | Not started; contract deferred |
| Real-workload testing | Not started |

### 8.1 Completed at this checkpoint

- r33/v9 architecture baseline and 118-requirement map.
- L-2 adopted; Decision A answered.
- G0 accepted; G1 fixture-contract passed.
- Production capture boundary traced and proposed (`PCB-01…42`).
- MEMORY-09 current-state audit and autonomous-classification plan.
- HSL requirement record (requirements only).
- Howl architecture/design package (proposed, not approved).
- Existing resource/VRAM/model-intelligence implementation surfaces identified.
- Consolidated readiness audit and reconciled MEMORY-09 decision register (`HFR-CP-2026-11-21-01`).

### 8.2 Active at this checkpoint

- This program direction and checkpoint.
- Next: Phase 1 work package defined (`HF1-PLAN-2026-11-21-01`) — state identity/context budget, followed by resource fault injection and a real-server smoke matrix.
- No implementation, G2, BGE or production integration authorization.

### 8.3 Blocked at this checkpoint

- G2 classification experiment — production capture boundary and corpus.
- Unresolved-class lifecycle — owner decision, L-2/Decision A amendment.
- BGE embedding experiment — separate model/corpus authorization.
- Production capture implementation — separate production authorization.
- Howl ↔ MEMORY-09 contract — HSL foundation and M5.
- Real-workload verification — H and I evidence.

### 8.4 Invariants to carry forward

**HPD-21 (Observed).**

- No class from incomplete capture.
- No fabricated knowledge to fill a missing field.
- No collapse of capture-incomplete into classification-unresolved.
- No authored-world eligibility from character speech subject matter.
- No inferred `told`; attribution is explicit or unresolved.
- Memory remains independent of inference state.
- No cloud dependency for ordinary operation.
- No unmeasured performance claim.
- No change to L-2, Decision A, r33, v9, the ledger or the map without the proper act.

### 8.5 Resume instruction

**HPD-22 (Calculated).** On return:

1. Read §2 and §8, then the readiness audit's §1, §10 and §13.
2. Compare the current tree and records to these checkpoints; do not assume they are still current.
3. Re-run the readiness audit's focused tests before starting new implementation.
4. If a decision is needed, record it; do not infer it from a subsystem's convenience.
5. If work has drifted, ask: does it reduce an unresolved decision or produce evidence for §6?

---

## 9. Immediate next move — consolidated readiness assessment

**HPD-23 (Owner direction).** Establish a consolidated **Howl foundation-readiness assessment** alongside the
**remaining MEMORY-09 decision register**.

**Output:** [`HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md`](HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md) — **completed 2026-11-21** at checkpoint `HFR-CP-2026-11-21-01`.

**Result:** A is not ready; B–D are partial/unverified; E is unestablished. No subsystem has end-to-end verification. The next evidence work is state identity + context budget, resource fault injection, and a real-server smoke matrix.

**Next package:** [`Howl Foundation Phase 1 — Chat State and Context Integrity`](HOWL-FOUNDATION-PHASE-1-CHAT-STATE-CONTEXT-INTEGRITY-2026-11-21.md) — definition complete; implementation not authorized.

### 9.1 Scope

| Section | Required content |
|---|---|
| **Readiness A — Chat/state** | HSL requirement-by-requirement status; current code evidence; branching/revert/retention gaps; pass/fail/unknown |
| **Readiness B — Multi-character** | Identity, sequencing, context flushing, voice separation, shared scenes, awareness, access boundaries; contamination risks |
| **Readiness C — Hardware** | Resource manager, telemetry, leases, reservations, OOM recovery; what is measured vs designed |
| **Readiness D — Inference** | KV cache, context sizing, sampling, stop sequences, batching, speculative decoding, quantization, backend support |
| **Readiness E — Performance** | Benchmarks that exist; benchmarks needed; no assumed gains |
| **MEMORY-09 register** | Every open decision, closed decision, blocked decision and superseded claim; dependency and authorization status |
| **Integration readiness** | Which §6 scenarios can be tested now, which need a foundation seam, which need M |

### 9.2 Evidence rules

**HPD-24 (Calculated).**

- Separate **Observed** code/runtime evidence from **Source-derived** design claims.
- A design package is not an implementation.
- A fixture is not production capture.
- A code path is not an acceptance result.
- A benchmark is not a writing-quality result.
- Unknown is a value; do not fill it with a plausible default.

### 9.3 Explicit non-goals for the next move

**HPD-25 (Observed).**

- No production code or API change.
- No schema or migration.
- No G2, BGE or model fetch.
- No corpus.
- No L-2/Decision A amendment.
- No claim that Howl is ready because a subsystem exists.
- No claim that MEMORY-09 is complete because the architecture is coherent.

---

## 10. Guardrails and stop conditions

**HPD-26 (Calculated).**

| Guardrail | Stop if… |
|---|---|
| **Destination** | The work optimizes an incomplete version of Howl instead of building toward the intended environment |
| **Uncertainty** | A missing capture field is being repaired by a classifier guess |
| **Attribution** | `observed`, `told`, `believed` or `authoritative` is being inferred from role or prose alone |
| **Authored status** | A chat message is becoming authored-world knowledge through subject matter |
| **State identity** | A branch, message or memory identity is derived from list position |
| **Lifecycle** | A chat-state transition is being inferred by MEMORY-09 from content, absence or ordering |
| **Hardware** | A performance or fit claim is made without measurement |
| **Scope** | Canonical files, L-2, Decision A or the 118 map are being changed without a separate authorized act |
| **Authorization** | Implementation, G2, BGE or production API work begins under this planning record |

---

## 11. Traceability and record status

**HPD-27 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Production code/schema edited | **0** |
| Requirements changed | **0** — map **118** |
| L-2 / Decision A amended | **0** |
| BGE / G2 started | **0** |
| Models fetched / corpus created | **0** |
| Tests executed | **0** — this is a direction/checkpoint record |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **program direction**, not canonical capture |

**HPD-28 (Calculated — limits).** The readiness matrix is record-based. It does not replace the §9 audit. The 2026-11-18
implementation figures are a dated audit, not a fresh measurement. The HSL gaps are from the 2026-11-19 requirement
record, not a fresh runtime test. No code was executed.

---

## 12. What “on track” looks like

**HPD-29 (Calculated — a practical test).** We are on track when:

1. The current checkpoint has a smaller set of unresolved **decisions**, not just more documents.
2. Each workstream can name its evidence: fixture, runtime, benchmark, or owner decision.
3. Integration scenarios are tested against real interfaces, not mocked around missing foundations.
4. Memory uncertainty survives contact with real conversations.
5. Hardware pressure does not silently corrupt memory or state.
6. A restart reconstructs the same authoritative state.
7. The system is being built toward the environment the owner wants to use.

**HPD-30 (Observed — the failure mode this record exists to prevent).** Repeatedly discovering that one unfinished
subsystem blocks another, after building an experiment around the limitation. The recovery index, the checkpoint and
the decision register exist to make that visible early.

---

**End of program direction.** The destination is one integrated local creative environment, not a collection of
promising subsystems (HPD-01). MEMORY-09 is design-strong but not executed; Howl has working pieces but no consolidated
readiness evidence (HPD-08, HPD-11). The work proceeds in parallel workstreams with explicit integration gates
(HPD-05…HPD-19). The next move is the consolidated readiness assessment and decision register, with no production
change, no G2, no BGE and no L-2/Decision A amendment (HPD-23…HPD-25). **Resume at `CP-2026-11-21-01`; search for
`HPD-CP-2026-11-21-01`.**
