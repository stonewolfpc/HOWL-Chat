# FC-5 — Capability, Resource, and Handoff Contract

## Contract paper `FC-5-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Contract paper — review draft |
| Authority | D6 owner decision; `HRA-DB-2026-11-21-01`; prior hardware design |
| Gate | M4 |
| Status | **PROPOSED — reservation questions unresolved** |
| Implementation authorization | **None** |
| Tests/models/process runs | **0** |

**Labels:** **Owner-approved** · **Proposed** · **Unresolved** · **Observed** · **Unestablished**.

**Identifier note:** Existing letter-suffixed clauses include `FC5-22a` and `FC5-22b`. This amendment adds `FC5-26a` as a related clause without renumbering any existing clause, so cross-references from other papers and the register remain valid. Clauses numbered above the highest existing clause (`FC5-27`–`FC5-31`) continue that sequence and are placed by subject rather than by number: numbering follows the highest existing clause while document order follows subject, so FC5-27–FC5-31 appear consecutively in §4 and FC5-26a appears in §10 after them despite its lower number. Identifiers of the form `FC5-CONFIG-RETENTION` identify questions, not clauses, and are cited as such.

---

## 1. Contract decision

**FC5-01 (Owner-approved — HRA-01; HRA decision summary).** Cross-modality resource arbitration and verified model handoff are foundational. Selection is intent; admission is a decision; loading is allocation; ready is an evidence claim; unloading is not release.

**FC5-02 (Owner-approved — D6).** The initial supported policy is **D6 option A**: one heavy-work instance; a second instance is read-only/light-work or refuses heavy admission. This is explicitly temporary until option B or option C passes recovery and concurrency acceptance. Those alternatives are a shared broker and an atomic lease registry; this paper does not assign either to a specific option letter. Option A is a temporary selection, not a permanent architecture.

**FC5-03 (Calculated — D6 scope).** D6 is an instance-level admission policy, not a one-resident-model rule. Supported multi-model residency remains required when resources and backend behavior safely allow it.

**FC5-04 (Calculated — D6 is not T0.5).** T0.5 proves process/job ownership: two instances do not kill each other's children. D6 does not prove that, and T0.5 does not prove global reservation accounting.

**FC5-05 (Calculated — D6 is not global accounting).** A per-instance ledger is not a global guarantee. Two instances can each believe the same free bytes exist unless they share an authority or explicitly limit concurrent heavy work.

---

## 2. Boundaries

**FC5-06 (Proposed — inherited from HRA).**

| Boundary | Owns | Must not own |
|---|---|---|
| Resource/safety governance | Admission, reservation, headroom, policy, cross-consumer arbitration | Model semantics, chat history |
| Model lifecycle adapter | Drain, unload, load, warmup, readiness, release evidence | Global admission |
| Conversation orchestration | Turn state, capability request, ordering, cancellation, work preservation | Hardware control |
| Modality services | Vision, audio, TTS, image, video, embeddings | Global resource policy |
| State/persistence | Conversation state, event history, recovery | Model choice |
| User control | Status, consequences, override, recovery | Hidden decisions |
| Evidence | Telemetry, receipts, release/readiness proof | Converting unknown to success |

**FC5-07 (Proposed).** No component may treat its own claim as proof from another component. A governor may approve a reservation without knowing how a runtime releases memory. A runtime may exit without proving the machine safe.

---

## 3. Reservation model

**FC5-08 (Proposed — logical vs physical).** A logical reservation is Howl's accounting commitment. A physical allocation is what the OS, driver, runtime and other applications actually do. A reservation never guarantees physical capacity.

**FC5-09 (Proposed — reservation record).**

| Field | Purpose |
|---|---|
| `reservation_id` | Stable identity |
| `owner` | Workload/request that owns the obligation |
| `capability` | Capability and version |
| `resource_envelope` | Weights, KV/context, workspace, device class |
| `state` | Requested/admitted/reserved/active/released/expired/unknown |
| `expiry` | Lease validity |
| `idempotency_key` | Duplicate/retry suppression |
| `evidence_ref` | Admission and release receipts |

**FC5-10 (Proposed).** A reservation is recorded atomically against a coherent snapshot. Within a single authoritative reservation ledger and its atomicity domain, the same capacity cannot be promised twice. This does not establish global coordination across independent instance ledgers; cross-instance safety remains governed by the active D6 policy and the unresolved reservation-authority questions. Allocation without a valid reservation is refused.

**FC5-11 (Proposed — physical evidence wins).** An unsafe observed state overrules a logical ledger that says capacity is free. If observation is unavailable, the ledger alone cannot prove safety.

---

## 4. Handoff phases

This section defines the handoff phase sequence and the two configuration records that accompany it: a record of the configuration decision taken for an operation, and a record of what Howl constructed for a launch attempt. The records are defined at the end of this section (FC5-27–FC5-31). They are not phases of the sequence in FC5-12, and no row of that table is satisfied by either record.

**FC5-12 (Proposed — required sequence).**

| Phase | Required outcome | Failure behavior |
|---|---|---|
| Intent | Capability, owner and work identified | Refuse invalid request |
| Assess | Envelope and confidence known | Mark unknown; apply conservative policy |
| Admit | Coexist/transition/queue/refuse decided | Never double-spend within the applicable reservation authority; cross-instance safety follows the active D6 policy |
| Reserve | Capacity claimed before allocation | No allocation without reservation |
| Drain | Outgoing work quiescent or explicitly interrupted | Keep resident; report busy |
| Release | Runtime resources released as far as observable | Mark partial/unknown |
| Verify release | Release proven or explicitly unknown | Do not claim freed capacity |
| Load | Artifact valid and runtime initialized | Explicit failed state |
| Warm/probe | Workload-specific readiness proven | Do not dispatch |
| Commit | Ready + valid reservation | Dispatch or refuse |
| Recover | Clean state, preserved work, explicit retry/refuse | No phantom ready |

**FC5-13 (Owner-approved — HRA-17).** Selection does not allocate. A selected model may be queued or unavailable; selection survives without residency.

**FC5-14 (Owner-approved — HRA-18).** Unload means actual release as far as the platform can demonstrate. Stopping inference is not release if weights, runtime allocations, caches, workspaces or job-owned processes remain.

**FC5-15 (Owner-approved — HRA-22).** Failed handoffs recover cleanly: no misleading ready state, no silent model/placement substitution, no destroyed user work.

**FC5-27 (Proposed — configuration decision record).**

An operation's configuration is decided in three distinguishable stages, and each stage is recorded:

| Stage | Recorded content |
|---|---|
| User intent | What the user asked for, captured at the operation boundary before application defaults are merged |
| Application defaults | The defaults Howl applied, each named with the identity of the policy that supplied it, distinguishable from user intent |
| Configuration grant | The configuration values Howl permits for this request, the adjustment made to each value where the grant diverges from the request, the stated reason for each adjustment, and the identity of the granting decision |

User intent is not reconstructed from merged or persisted settings. A stored settings value that already contains an application default is not evidence that the user chose that value. Where intent was not captured — for example because the operation began from persisted settings rather than from a request — the record states that, and does not present a merged value as intent.

The configuration grant is not FC-5 admission. Admission decides whether a workload may proceed under capability and resource policy. The configuration grant decides which configuration values Howl permits for that request. The two decisions may constrain one another, but neither substitutes for the other and neither is evidence for the other. A grant is also not a reservation state: it is not recorded in FC5-09's `state` values, and an existing reservation implies no grant.

The configuration grant is not a claim-ladder level (FC5-22a). It is a configuration decision and carries no evidence maturity of its own.

This clause defines the record's content and the distinction between its stages. It does not specify which component writes the record, where it is stored, or in what format beyond the distinctions above (FC5-30, FC5-31).

---

**FC5-28 (Proposed — launch construction record).**

For each launch attempt, Howl records what it constructed, expressed per relevant setting or argument:

| Disposition | Required content |
|---|---|
| `emitted` | The option and the value emitted |
| `omitted` | The option and an explicit reason |
| `transformed` | The option, the requested value, the emitted value, the kind of transformation, and its provenance |

The kind of transformation distinguishes a semantic value change from a representation-only change. A semantic change alters the effective value or its meaning. A representation-only change expresses the same effective value differently — spelling, ordering, formatting, or the field or flag form used — and does not alter what the runtime would be asked to do. The two are recorded as different kinds, because they carry different consequences for a later reader.

An omission reason preserves the distinction that explains it. At minimum the record separates: never requested; requested, not emitted, and left at the application default; conditionally excluded by a stated condition; and requested but dropped as unsupported. These are distinct reasons and are not collapsed into one unspecified reason.

The second of these is an omission only while the default remains in effect and no value is emitted in place of the request. Where Howl substitutes an application default for the requested value and emits a value, the outcome is `transformed`, not `omitted`. Its kind is determined by comparing the effective values, not by the fact of substitution: where the substituted default differs from the requested value in effect, the change is semantic; where the substituted default is equivalent to the requested value, the substitution is still recorded, but the kind is representation-only, because the effective configuration did not change. A substitution is never described as a semantic change on the strength of the substitution alone.

Where construction performs a derivation — applying a default, clamping a value, resolving an unsupported request, conditionally excluding an option — the record states that derivation's provenance: the stage that performed it, the identity of the governing policy where one exists, and the build identity associated with the construction record, reusing FC-6's `build_identity` with the meaning FC6-06 and FC6-07 already give it. Where no governing policy identity exists, the record states that rather than inventing one.

The build identity recorded here is the one associated with the record. Recording it does not assert that this build performed the derivation for this attempt. Where that relationship is established, the record states it at the strength the evidence supports; where it is not, the record states the association only and its provenance is correspondingly weaker. This clause creates no identity for a derivation and no registry of derivation rules.

Record B describes construction only. What that construction does not establish is stated in FC5-30.

---

**FC5-29 (Proposed — identities, references and retry).**

Each configuration decision record has its own identity. Each launch construction record has its own identity and references the configuration decision record it derives from.

The construction identity is FC-5's own and is capability-agnostic: it applies to text, vision, audio and other operations alike, and there is one identity per construction attempt. It is not `request_id`, `prompt_plan_id` or `reservation_id`, which identify a prompt-dispatch attempt (FC2-05b), a composition artifact and a reservation respectively; none of those scopes is a launch construction attempt.

A reservation reference is recorded when a reservation has been established for the attempt. It is never universally mandatory. FC5-12 places reservation before load and does not state where construction falls relative to it, so a construction record is not invalid for want of a reservation, and an absent reference is not a claim that no reservation existed.

One configuration decision can support more than one construction attempt. The configuration decision identity remains the common reference across those attempts; each construction record retains its own distinct identity. This identity relationship does not establish any retention or persistence guarantee (FC5-31).

A genuine change to the configuration decision produces a new decision identity, and the earlier decision record is not rewritten to represent it. A failed construction record is subject to whatever retention policy is eventually established (FC5-31). This clause does not establish that policy.

---

**FC5-30 (Proposed — configuration records are not runtime evidence).**

Neither record authorizes, implements or proves any of the following, and neither may be cited as evidence for any of them:

- configuration recording or storage infrastructure;
- observation or monitoring of a process;
- collection of runtime telemetry;
- verification that an argument, option or value was accepted by a runtime;
- backend or execution identity;
- allocation, loading, hardware execution, or workload readiness.

A configuration grant is a decision about permitted values. A launch construction record is a description of what Howl constructed. Neither establishes what a runtime accepted, what was allocated, what loaded, what executed, or whether a workload became ready. Those runtime facts remain governed by their applicable contracts. This clause attributes no evidence obligation to any other contract or phase, and does not claim that any phase or paper owes backend identity, hardware-execution evidence or runtime telemetry.

Readiness remains FC-5's own phase claim (FC5-22). This clause does not extend FC-5 to backend identity, hardware execution or runtime telemetry, and neither record becomes readiness evidence because FC-5 owns readiness. The records are not promoted into readiness claims, not even when they are complete and mutually consistent.

---

**FC5-31 (Unresolved — configuration-record retention).**

This clause defines no retention behaviour and does not qualify FC5-27–FC5-30. Those clauses require the records to be defined and to carry their identities and references; they do not establish that a record is written to durable storage, that it survives a process exit or restart, that it is retained for any duration, where it is held, what owns it, what scope it has, or how it is redacted.

`DQ-HRA-04` asks which **reservation** records survive a crash; it does not answer these questions for configuration records. D2's retention-by-default governs authored content and does not reach application settings or launch records. Configuration-record retention is therefore not established by the provisions governing these records: nothing in FC5-27–FC5-30, `DQ-HRA-04` or D2's retention-by-default requires these records to persist. This clause does not decide whether any other provision constrains where or how they may be stored.

Open question, minted here: `FC5-CONFIG-RETENTION` — what retention, storage, scope and redaction policy applies to configuration decision records and launch construction records?

This clause is not a `DQ-HRA-*` question, is not inherited from `HRA-DB-2026-11-21-01`, and decides none of the points it lists.

---

## 5. D6 temporary policy and exit

**FC5-16 (Owner-approved — D6).** While D6 option A is in force:

- one instance is the supported heavy-work instance;
- a second instance may run light/read-only work or refuse heavy admission;
- the refusal is explicit and user-visible;
- no instance claims global reservation safety from its own ledger.

**FC5-17 (Owner-approved — D6 exit).** Options B and C become eligible only after recovery and concurrency acceptance. Their definitions remain with the D6 intake decision; this paper does not select either. Implementation completion, a successful demo or an apparently working broker is not evidence.

**FC5-18 (Proposed — multi-model residency).** Under D6 option A, multiple resident models remain permitted when the supported instance can prove safe coexistence. D6 is not a licence to serialize all model residency permanently.

**FC5-19 (Proposed — T0.5 independence).** T0.5 acceptance requires two independent instances, each with an owned child or approved owned test child. A second instance refusing heavy work has not exercised the ownership boundary and cannot close T0.5.

---

## 6. Capability envelope

**FC5-20 (Proposed).** A capability declares enough for admission without hardcoding the modality into chat:

- identity/version;
- inputs/outputs/provenance;
- weights/runtime/context/workspace/storage/thermal envelope;
- placement support;
- lifecycle and readiness probe;
- concurrency/interruptibility;
- priority/fairness;
- privacy/local processing;
- failure modes;
- user-work preservation.

**FC5-21 (Proposed).** A capability that cannot declare its envelope is not ready for admission. It may run only under an explicit conservative policy that says so.

---

## 7. Evidence ladder

**FC5-22 (Owner-approved — HRA-34).** “Ready” is an evidence claim, not a status label. The same applies to persistence, modality readiness, state transitions and memory writes.

**FC5-22a (Proposed — evidence ladder).**

| Level | Permitted claim |
|---|---|
| Requested | The user asked |
| Admitted | Howl decided to try |
| Attempted | The operation started |
| Observed | Something appears to have happened |
| Verified | A workload-appropriate check confirms the outcome |
| Unknown | Evidence is missing, insufficient or contradictory; the outcome is not established |
| Failed | Evidence establishes that the operation did not succeed |

Neither `unknown` nor `failed` may be promoted to ready, complete or successful.

**FC5-22b (Unresolved — memory boundary).** Naming memory writes in FC5-22 establishes no memory-record ownership, authority, lifecycle, applicability or integration requirement, and does not preempt FC-7 or the delegated D8 decision. Evidence requirements for independently authorized memory operations are unaffected.

**FC5-23 (Proposed).** No promotion by adjacency: process start is not readiness; unload return is not release; UI update is not durable commit; a log line is not proof.

---

## 8. Unresolved reservation questions

**FC5-24 (Unresolved — inherited from HRA).**

| ID | Question | Status |
|---|---|---|
| `DQ-HRA-01` | Which component is the sole writer of reservation state? | Unresolved |
| `DQ-HRA-02` | Can operations share a parent reservation, and under what sub-budget rules? | Unresolved |
| `DQ-HRA-03` | What is the required response to an estimate overrun? | Unresolved |
| `DQ-HRA-04` | Which reservation records survive a crash, and how are they reconciled? | Unresolved |
| `DQ-HRA-05` | What evidence distinguishes a stale reservation from an orphaned allocation? | Unresolved |
| `DQ-HRA-06` | How does Howl respond when another application takes memory after admission? | Unresolved |
| `DQ-HRA-07` | Is multi-instance concurrent heavy work supported, and through which authority? | Unresolved |
| `DQ-HRA-08` | Which readiness probe is valid for each capability and workload? | Unresolved |
| `DQ-HRA-09` | Which release evidence is obtainable on Windows per runtime/driver? | Unresolved |
| `DQ-HRA-10` | What are the non-overridable safety limits, as distinct from warned soft limits? | Unresolved |
| `DQ-HRA-11` | What is the user-work snapshot/restore contract across a handoff? | Unresolved |
| `DQ-HRA-12` | How are future modality capability envelopes declared and versioned? | Unresolved |

**FC5-25 (Calculated).** D6 option A is a temporary answer to one part of `DQ-HRA-07`; it does not close the question, and it does not close `DQ-HRA-02` through `DQ-HRA-06`.

---

## 9. Acceptance design — proposed, not run

1. Selecting a model allocates nothing before admission.
2. Two concurrent requests cannot reserve the same free bytes.
3. Coexistence is allowed when evidence supports it; otherwise queue/transition/refuse.
4. A busy capability is drained or refused, never silently killed.
5. Unload reports verified/partial/unknown; unknown is never shown as free.
6. A process start does not pass readiness.
7. Failed load/readiness preserves selection and user work; no silent substitution.
8. Missing telemetry produces a defined conservative behavior.
9. Soft limits are overridable with consequences; hard limits are not.
10. Two instances either coordinate or explicitly limit heavy work.
11. An orphan is reported and recoverable, not killed because unexplained.
12. T0.5 evidence is separate from reservation accounting.

---

## 10. Traceability and non-authorization

| Link | Artifact |
|---|---|
| Owner decision | D6 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Resource baseline | `HRA-DB-2026-11-21-01` |
| Prior architecture | `AMENDMENT-01`; `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` |
| T0.5 acceptance | `HCDM-M0A-AR-2026-11-21-01` |
| Program map | `HCDM-2026-11-21-01` |
| Evidence | FC-6 (this set) |

**FC5-26 (Observed — no authorization).** This paper does not implement a broker, lease registry, process supervisor, telemetry provider, model loader or admission gate. It does not resolve the DQ-HRA questions, authorize M4, or claim T0.5 closure.

**FC5-26a (Proposed — boundary pointer).** FC5-26's exclusion list is not exhaustive. The configuration records defined in §4 introduce boundaries that list does not name: FC5-30 states what those records neither authorize nor prove, and it applies to both records. This clause points to that boundary and does not restate it.

**End of FC-5.** Admission is not residency, readiness is not process start, release is not unload return, and a per-instance ledger is not global safety.
