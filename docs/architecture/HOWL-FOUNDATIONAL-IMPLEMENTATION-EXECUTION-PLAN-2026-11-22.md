# Howl — Foundational Implementation Execution Plan

## Plan `HFI-EP-2026-11-22-01`

| Field | Value |
|---|---|
| Date | **2026-11-22** |
| Type | Execution plan — planning only |
| Supersedes as the active activity | Chronological HEKI discovery (paused at `HEKI-CP-2026-11-22-03`) |
| Basis | `HFR-CP-2026-11-21-01`, `HFS-AUDIT-2026-11-21-01`, ADR-001, ADR-002, FC-1…FC-6, `HCDM-2026-11-21-01`, `HRA-DB-2026-11-21-01`, `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Implementation authorization | **None.** This plan authorizes no code, schema, API, migration, test run, build, launch or contract amendment |
| First repository-changing task | Defined here as `TASK-01` — **prepared, not authorized**. Requires explicit owner approval |
| Recovery token | **`HFI-EP-2026-11-22-01`** — search the repository for this token |

**One sentence.** The architecture is decided; the ownership is not yet enforceable; so the next work is to make
ownership executable, one bounded slice at a time, with each slice verified before the next begins.

---

## 1. Why the activity changes

Two workstreams replace one. This is not a reduction in the value of the recovery work; it is a change in which
workstream is primary.

| Workstream | Goal | Cadence |
|---|---|---|
| **Primary — Howl foundational architecture** | Turn decided architecture into bounded, tested, owned implementation | Continuous; each slice ends in a handoff |
| **Supporting — HEKI / MEMORY-09 recovery** | Preserve research, locate evidence when implementation asks, record findings | On demand; no obligation to finish the corpus |

**Stopping rule.** Chronological discovery stops at `HEKI-CP-2026-11-22-03` (22 of 808 verified). Tranche 04
(62 MEMORY-09 records dated 2026-11-20) is **deferred, not cancelled**. Nothing is discarded, no identity is
renumbered, and the ledger stays valid as a census.

**What HEKI stops doing:** it no longer gates implementation. `README.md` R7 ("do not begin implementation on an
incomplete reconstruction") is **superseded as a gate** and retained as a principle: an implementation task must
state which evidence it relies on, and may proceed when that evidence — not the whole corpus — is present.

**What HEKI keeps doing:** answering "what does the record actually say" for a specific question, with the
citation and the status distinction intact. Its two proven capabilities are exactly what the primary workstream
needs: catching a document that contradicts the state of a file (§`ERRATA-CANDIDATES.md` E-01), and finding
reachability that a link graph misses (E-02, T3-F1).

---

## 2. Basis, and the status of each input

Nothing below is promoted by being listed. A `Proposed` paper stays proposed.

| Input | Identity | Status as written | What the plan takes from it |
|---|---|---|---|
| Readiness audit | `HFR-CP-2026-11-21-01` | Audit complete; **not ready** for integrated real-workload acceptance | A: history E2, state machine **E0**; B partial; C partial, `internal/resource` has no tests; D partial; E **E0**; production changes **0** |
| Original-scope audit | `HFS-AUDIT-2026-11-21-01` | Read-only audit | The **27-item** T0–T3 union; T0.5 preserved; V2's omission is not closure |
| Backend rebuild | **ADR-001** | Accepted architecture; **implementation not authorized** | 5 roots, 6 services + 2 primitives, contracts A–F, D1–D5, L0–L3 identity |
| Storage/navigation | **ADR-002** | Accepted policy; **implementation not authorized** | 3 state families, N1–N19, adoption records, no cross-writes |
| Contract papers | FC-1…FC-6 | **PROPOSED — not approved** (6 papers, 0 approvals) | FC-1 ledger authority; FC-2 PromptPlan; FC-3 isolation; FC-4 canonical library; FC-5 admission vocabulary; FC-6 evidence ladder |
| Contract register | `HF-CONTRACTS-2026-11-21-01` | Review-complete; none approved; FC-7 absent (D8) | Dependency order and the deliberate hole |
| Decision register | `HF-CONTRACT-INTAKE-2026-11-20-01` | Owner decisions recorded; contracts unapproved; M0A unauthorized | D1, D2 (as narrowed), D3, D5, D6 selected; D4 direction only; D7 clarified; D8 delegated |
| Delivery map | `HCDM-2026-11-21-01` | Planning map | M0A–M8 ordering and gate names |
| Resource baseline | `HRA-DB-2026-11-21-01` | Discovery baseline; 12 open questions | `DQ-HRA-*`; D6 is not T0.5 and not global accounting |

**Unresolved by design, carried forward:** the six contract papers are unapproved. Implementation therefore
proceeds against **ADR-001/ADR-002** (accepted) and treats FC-1…FC-6 as the contract the implementation must
*satisfy*, not as authority it can cite. Where a paper is `Proposed`, the implementation records which clause it
implements and does not claim the clause is adopted.

---

## 3. Ownership contracts — the unit of work

Every subsystem gets one contract in the table in §4. Seven questions, no blanks, no "shared".

| Question | Required answer | Failure if unanswered |
|---|---|---|
| **Owner** | Exactly one authoritative component | Two writers, or none |
| **State** | The information it owns, named | Ownership by vague domain |
| **Writers** | Who may change it, by what path | Silent bypass |
| **Readers** | Who consumes it, and how | Uncontrolled coupling |
| **Lifecycle** | Who creates, updates, destroys | Orphans and ghosts |
| **Boundaries** | What it must **not** own | Responsibility creep |
| **Failure** | Behaviour when the owner cannot complete | Silent degradation or a phantom success |

**Principle:** one authoritative owner per responsibility; delegation is explicit and scoped; **no silent shared
ownership**. A second implementation of a responsibility is a defect to be removed, not a fallback to be kept.

**Why this comes before code.** ADR-001 §1 states the finding plainly: the defects "cluster where one
responsibility has more than one implementation." Writing features before fixing that reproduces it.

**Existing owners are reused, not re-invented.** The inventory builds on ADR-001 §4 (service table with `Owns` /
`Must not own` / `Depends on`), ADR-002 §2 (three state families), FC-1 (ledger vs library vs snapshot) and FC-4
(canonical write authority). This plan adds no seventh service and no competing architecture document.

---

## 4. Ownership inventory

Status: **U — unimplemented (decided, no code)** · **P — partial implementation exists** · **C — conflict, two
candidates** · **X — no owner exists anywhere**.

| # | Subsystem | Owner (decided) | State it owns | Writers | Readers | Failure behaviour | Status |
|---|---|---|---|---|---|---|---|
| O-01 | Path resolution | **PathService** | Root set, mode, category→location, class labels, containment checks | none (resolution is pure) | every store and service | typed failure naming category + remedy; never relocates | U |
| O-02 | Configuration (app prefs) | **ConfigAuthority** | App-preferred settings, schema, defaults, reset semantics, **authority record**, navigation state | ConfigAuthority only, via DurableStore | RuntimeResolver, UI | `Saved`/`NotSaved`/`RefusedUnknownSchema`; no boolean | C |
| O-03 | Model/preset config | **ModelConfigService** | Presets (stable IDs), model↔preset associations, per-model custom inputs | ModelConfigService | RuntimeResolver, UI | unmatched preset preserved and marked; never auto-rebound | C |
| O-04 | Machine config | **MachineConfigService** | Engine paths, ports, machine-scoped desired values | MachineConfigService | EngineLifecycle, RuntimeResolver | validation failure refuses save; **never restarts** as a side effect | P |
| O-05 | Runtime resolution | **RuntimeResolver** | Resolution pipeline + transformation and launch **records** (not intent) | appends records only | diagnostics, evidence | each stage recorded separately; `unestablished` is terminal | U |
| O-06 | Engine lifecycle | **EngineLifecycle** (coordinator) | Start/stop/restart and outcome reporting | explicit call only | UI, diagnostics | explicit failed state; no phantom ready | P |
| O-07 | Document durability | **DurableStore** (primitive) | Atomicity, backup, checksum, version envelope, single-writer lock | any store, under the lock | any store | `Busy` on contention, **no bypass**; corrupted → recovered or refused | U |
| O-08 | Containment | **ContainmentChecker** (primitive) | Reparse inspection, canonical identity pinning | none | ExternalDataGuard, PathService | fails closed | P |
| O-09 | External writes | **ExternalDataGuard** | Operation-scoped grants, consent records, created-file manifest | grant issuance only | importers, installer | expired/absent grant refuses and records the refusal | U |
| O-10 | Conversation state | **Conversation authority** (FC-1) | Append-only event ledger; derived current selection | append-only commit path | composition, restore, UI | append failure = no commit; no in-place edit | **C** |
| O-11 | Authored content | **Canonical library** (FC-4) | Authored objects + immutable revisions | library write path only | composition, snapshots | `needs-review` on identity change; no silent rebind | **C** |
| O-12 | Context composition | **Composition service** (FC-2) | PromptPlan for a turn, bound to exact revision IDs + tokenizer identity | composition service | dispatch, evidence | **fails closed** when the tokenizer is unavailable — `unavailable`, not an estimate | **X** |
| O-13 | Speaker isolation | **Per-speaker execution context** (FC-3) | Private scope, projected transcript, token allocation per speaker | orchestrator per speaker | dispatch | hard flush on switch; no cross-speaker reuse without evidence | **X** |
| O-14 | Admission/reservation | **Unresolved** (`DQ-HRA-01`) | Reservations, admission decisions | — | — | D6 option A is the temporary policy; a second instance refuses heavy work **explicitly** | **C** |
| O-15 | Memory records | **Unresolved** (`FC4-04a`, FC-7 blocked by D8) | — | — | — | no authority is granted or denied; MEMORY-09 stays separate | **X** |
| O-16 | Evidence/receipts | **Unassigned** | Acceptance receipts, claim levels, launch records | whoever produces an observation | audit, acceptance | a missing receipt is `Unknown`, never success | **X** |

**Not in doubt, and therefore not in the table:** `worldstore`, `contentimport`, `lorebookstore`, `library`
retain their existing ownership model and journals (ADR-001 §4). The image-engine process lifecycle, llama.cpp
integration, audio settings contract and error/diagnostics surfaces are **retained unchanged**.

---

## 5. Conflicts and missing owners

Each is stated with its evidence and a disposition. None is resolved by this plan.

**C-1 · Character is claimed by two authorities.** FC4-04 places Character in the canonical library.
ADR-001 §11.2 says preserve the existing root in place and *"do not move, merge into the library system, or
overwrite"*, with the permanent root `UN`. **Disposition:** Phase 1 treats characters as **read-only references**
and implements neither claim. The conflict must be settled before any character write path is built (M5, not
Phase 1).

**C-2 · Conversation state has five candidate implementations and no owner.** Today: in-memory `a.chatMessages`,
a second context in `chat.Service`, the `history.json` snapshot, the `history.ndjson` append log, and frontend
DOM/localStorage paths. FC-1 says the ledger is durable truth and every projection is derived. **Disposition:**
Phase 1B creates the ledger owner and the projections become explicitly derived. Removing the other four is part
of the slice, not a follow-up.

**C-3 · Reservation authority.** `internal/resource.Manager` tracks registrations; `app-vram-arbitration.go`
performs one LLM/image handoff. Neither is a reservation authority, and `DQ-HRA-01` is open. **Disposition:**
Phase 1 requires an **admission interface**, not a reservation authority. Where a turn needs a model, Phase 1
records the request and the refusal path; it does not claim reservation safety.

**C-4 · Configuration authority is four-way overlapping with temporal precedence.** ADR-001 §1 names them:
unified settings, legacy model/sampler files, frontend `localStorage`, live in-memory maps. **Disposition:**
Phase 1 does not migrate configuration. It implements PathService + DurableStore and leaves the migration
sequence (ADR-001 §12, S0–S4) as its own slice, because a half-migrated configuration is worse than a legacy one.

**M-1 · No owner for context composition.** FC-2 is `Proposed`; no component owns the PromptPlan. Current code
estimates by characters-per-token and the backend has placeholder trimming. **Disposition:** Phase 1C assigns
the owner. The tokenizer capability must return an explicit `unavailable` — inventing an estimate is prohibited.

**M-2 · No owner for speaker isolation.** FC-3 is `Proposed`; the inspected backend send path selects
`characterNames[0]` and there is no per-speaker context proof. **Disposition:** Phase 1C. Multi-character
orchestration stays untested until this exists.

**M-3 · No owner for evidence/receipts.** FC-6 defines claim levels and receipt schemas, but no store owns them.
ADR-001 §9 places launch records in R2. **Disposition:** Phase 1 assigns receipts to a store under R2 and records
only what it observes.

**M-4 · Memory-record authority is deliberately absent.** `FC4-04a` grants and denies nothing; FC-7 waits on D8.
**Disposition:** unchanged. Implementation must not fill this hole by convenience, and the existing lorebook is
**not** the MEMORY-09 architecture.

---

## 6. Delivery sequence

Five vertical slices. Each ends in a handoff (§12) and is verified before the next begins.

| # | Slice | What it establishes | Gate | Why here |
|---|---|---|---|---|
| 1 | **Chat state and context integrity** | Identity, state ownership, transitions, isolation, tokenizer boundary | M1/M2/M3 | Everything downstream composes from it; today's state machine is `E0` |
| 2 | **Resource arbitration and hardware awareness** | VRAM/RAM awareness, resource ownership, model lifecycle, safe failure | M4 | Needs a state owner to hand work to; today C is partial and untested |
| 3 | **Inference and runtime configuration** | Family resolution, prompt construction, sampling config, KV-cache, verified outcomes | M2/M6 | Needs admission; today stages 7–9 are `UN` |
| 4 | **Multi-character orchestration** | Queues, turn ownership, flushing, concurrency, consistency | M3 | Needs isolation + admission; today untested |
| 5 | **Integration and broader foundation** | Coupled workflows, remaining T2/T3 items, M0A/M0B truth | M0A/M0B/M5 | Only after the pieces are independently verified |

**Proposed, not asserted.** This order follows the dependency direction in ADR-001 §4 and the gate map in HCDM.
Where the dependency map disagrees, the map wins and this table is corrected — an ordering claim is not an
architecture claim.

---

## 7. Phase 1 — the bounded slice

**Objective:** a conversation whose state is owned, identified, persisted, recoverable and composed with a real
tokenizer boundary, for a single speaker and for multiple speakers without leakage.

**In scope**

1. Stable identity for conversations, characters-as-references, turns and states (FC1-05); no positional identity
   (FC1-06).
2. Append-only ledger as durable truth; current selection as a projection (FC1-01, FC1-13a).
3. Authoritative transitions with explicit preconditions and recorded outcomes (FC1-14).
4. Speaker switch flushes context; no cross-speaker context without adapter evidence (FC-3).
5. Context construction with the actual model tokenizer; `unavailable` when it is not present.
6. Persistence and recovery of exactly the state above, with a stated durability level (ADR-001 D1/D3).
7. Failure paths: interrupted commit, interrupted transition, corrupt snapshot plus valid log, duplicate
   submission, mid-stream exit.

**Out of scope, explicitly**

Authoring/character write paths (C-1) · configuration migration (§5 C-4) · reservation accounting (C-3) ·
memory-record authority (M-4) · portable mode · modality work · any UI redesign · any contract amendment.

**The slice splits into three tasks, because two prerequisites are not optional.**

| Task | Content | Depends on |
|---|---|---|
| **1A** | PathService + DurableStore (ADR-001 contracts A, B, D) | — |
| **1B** | Conversation identity + ledger + projections (FC-1) | 1A |
| **1C** | Composition, tokenizer boundary, speaker isolation (FC-2, FC-3) | 1B |

**Why 1A is first, and this is a deviation from the stated first target.** Stable identity and recovery both
require a place to write and a way to prove the write is atomic and single-writer. Without 1A, every Phase 1
acceptance test would have to be demonstrated against the very paths ADR-001 declares untrustworthy. 1A is small,
fully specified, and has **no unresolved blockers** — which is what makes it a safe first slice. This is
**flagged for owner approval**, because it changes the first implemented target.

---

## 8. Minimum dependency chain for Phase 1

```
PathService (O-01) ──► DurableStore (O-07) ──► Conversation authority (O-10)
                                                    │
                        Canonical library refs (O-11, read-only in Phase 1)
                                                    │
                                    Composition (O-12) ──► Speaker isolation (O-13)
                                                    │
                        Receipts (O-16)  ◄──────────┘
```

Not required for Phase 1: EngineLifecycle (O-06), ExternalDataGuard (O-09) beyond compile-time availability,
ModelConfigService (O-03), MachineConfigService (O-04), RuntimeResolver (O-05) beyond recording, admission
authority (O-14) beyond an interface.

---

## 9. `TASK-01` — PathService substrate

**Prepared, not authorized.** No repository change occurs until the owner approves this task boundary.

**Goal.** One resolved storage identity, one root set per mode, one containment rule — with no migration, no
relocation, and no change to existing behaviour.

**In scope.** Mode detection (Installed / Isolated) · root resolution for R1–R5 · category→location through the
Windows known-folder API · class labelling · containment check by `(volume serial, file ID)` · typed failures ·
diagnostics output of `(stateRoot, appVersion, schemaVersion)`.

**Out of scope.** DurableStore (that is 1A's second half) · migration S0–S4 · adoption records for chats and
characters · any picker or navigation change · deleting or moving any file · the `chat` package.

**Acceptance criteria — observable**

**Expanded and superseded** by the owner's safeguard B (containment, external-path handling, failure behaviour and
compatibility, not merely successful resolution). The authoritative table, with measured results, is
`TASK-01-PATHSERVICE-OWNERSHIP-AND-COMPLETION-2026-11-22.md` §4. The eight criteria below remain as the original
statement of intent.

**Verification evidence required.** `go test` on the new package with the table above as named cases; the
three-directory diagnostic comparison; the before/after directory snapshots as files; a short record of what was
**not** verified (no model, no engine, no UI, no network).

**Stop conditions.** If the toolchain cannot resolve known folders or file IDs from Go, **stop and report** — do
not fall back to `%USERPROFILE%` string joins, and do not degrade containment to string comparison.

---

## 10. What actually blocks what

Only questions that stop the named task. Everything else waits.

| Question | Blocks | Status | Narrow investigation? |
|---|---|---|---|
| Any ADR-001 open item | `TASK-01` | **None block it** | No — 1A is fully specified |
| `FC1-12` ordering implementation (sequential assignment, writer authority) | `TASK-02` (ledger) | Unresolved in FC-1 | Yes — bounded; needed before ledger design |
| `FC1-07` legacy identity | `TASK-02` migration only | Unresolved | No — legacy stays unresolved, read-only |
| `FC1-SELECTION-AFTER-DELETION` | `TASK-02` selection path | Unresolved | Yes — one decision, one state |
| `FC1-14b`/`FC1-14c` branch and turn relationships | `TASK-02` branching | Unresolved | No — defer branching out of Phase 1 |
| `DQ-X01` tokenizer authority (D4) | `TASK-03` real budget proof | Owner direction only | Yes — but **not** the `unavailable` boundary, which is buildable now |
| `FC2` technology selection | `TASK-03` | `Proposed` | No — fail-closed first, select later |
| `DQ-HRA-01` reservation writer | Slice 2 | Unresolved | No — not needed for Phase 1 |
| `FC4-17` bypass inventory | M0B | **Unestablished** | Yes — but not Phase 1 |

**The important result:** the first slice has **zero blocking questions**. The plan chooses 1A partly because
that is true, and because a first task that is blocked on an owner decision is not a task.

---

## 11. Missing-information protocol

When implementation meets something unresolved:

1. **Stop at the specific uncertainty.** Not the surrounding topic — the one question.
2. **Name the requirement, component and ownership boundary** it touches.
3. **Retrieve the evidence** from HEKI (ledger, mentions, filerefs, `RELATIONS.md`) or from the MEMORY-09 records.
4. **If evidence is insufficient, run a narrowly scoped investigation** — one question, stated before it starts.
5. **Record the result and its provenance**, including what it did not establish.
6. **Resume the task.**

Prohibited: reopening the research programme for one question; inventing an answer to keep coding; and promoting
an investigation's result to a decision without the owner.

---

## 12. Durable handoff — required for every task

Five fields, no exceptions. A context reset should cost minutes.

| Field | Content |
|---|---|
| **Starting point** | Exact recovery token, the architecture documents relied on, and the verified state of the tree |
| **Ownership** | Owner, state owned, writers, readers, lifecycle, what it must **not** own, failure behaviour |
| **Acceptance criteria** | Observable behaviour with the test that demonstrates it — never "done" |
| **Evidence** | Test output, measurements, before/after artifacts, and **known limitations** |
| **Next action** | The exact next task, unresolved blockers, and how to resume without re-discovery |

**This plan's own recovery point:**

> Resume by reading this file's §4 (ownership), §7 (slice), §9 (`TASK-01`) and §10 (blockers). Then read
> `ADR-001` §4–§5 for 1A's specification. The last completed work is `HEKI-CP-2026-11-22-03`; HEKI is in
> supporting mode. Nothing is implemented yet: the tree matches the census (0 of 808 files differ).

---

## 13. Constraints carried forward from the owner's direction

- **T0.5 stays in the 27-item scope** (`HFS-AUDIT`). It is not closed by M0A's absence.
- **T2.2 remains open**, including its 18 ghost IDs and 30 references *(owner-stated 2026-11-22; not re-verified
  in this pass — recorded as `Unestablished` here and verified when M0B begins)*.
- **Character-based budgeting is not tokenizer-based budgeting.** The current estimate is not a measurement.
- **Multi-character orchestration is untested** until slice 4 produces evidence.
- **The lorebook is not the MEMORY-09 architecture.**
- **MEMORY-09 integration is not authorized** by rebuilding the foundation.
- **The six contract papers remain unapproved.** Implementing their clauses does not approve them.

---

## 14. Authorization

This plan authorizes **nothing executable**. It defines the workstreams, the ownership contracts, the slice, and
the first task boundary.

**The single next owner decision:** approve or amend `TASK-01` (§9) as the first repository-changing task. On
approval, the work proceeds as 1A → 1B → 1C, each with its own handoff, and slice 2 does not begin until Phase 1's
acceptance evidence exists.

FC-5 remains `PROPOSED`; `FC5-31` remains `Unresolved`; R6's text is applied but the paper is unapproved; ADR-001
and ADR-002 remain accepted architecture with implementation unauthorized. Nothing in this plan changes any of
those statuses.

---

## Task status (appended 2026-11-22)

| Task | State | Evidence | Next |
|---|---|---|---|
| `TASK-01` PathService | **Closed** — accepted after audit | `TASK-01-AUDIT-2026-11-22.md` §8; `HEKI-CP-2026-11-22-08` | none |
| `TASK-02` DurableStore | **Primitive complete, not adopted** | `TASK-02-READER-SHARING-DECISION-2026-11-22.md`, `TASK-02-DURABLESTORE-2026-11-22.md`; `HEKI-CP-2026-11-22-09` | adoption at `internal/settings` — **not authorized** |

Two things this table deliberately does not claim: that `TASK-02` changed application behaviour (no call site uses
the store), and that power-loss durability is established (`UN`).

Census distinction unchanged: historical census **808**, current corpus **813**, reserved placeholders
`HEKI-D-0809`–`0813` are not ledger nodes and have no recorded hashes. **813 is not the census.**

---

## Bookkeeping corrections (appended 2026-11-22, owner-approved)

No earlier text in this plan was rewritten. Each correction names what it supersedes, so the earlier statement stays
readable next to its replacement.

| # | Correction | Supersedes |
|---|---|---|
| 1 | **Task numbering.** The ledger — "conversation identity + ledger + projections" (§7, task 1B) — is **`TASK-03`**. | §10's row "`FC1-12` ordering implementation … `TASK-02` (ledger)", and its three neighbours naming `TASK-02` for the ledger. `TASK-02` is DurableStore and is closed |
| 2 | Consequently, §10's "`TASK-03` real budget proof" is **`TASK-04`**. | The `DQ-X01` row of §10 |
| 3 | **`TASK-02` status is closed**, not "primitive complete, not adopted". | The `TASK-02` row of the appended "Task status" table |
| 4 | **Corpus is 814** documents (808 census + 6 unindexed), not 813. Reserved placeholders are `HEKI-D-0809`–`0814`. | The census sentence of the appended "Task status" table |
| 5 | §12's "Nothing is implemented yet: the tree matches the census (0 of 808 files differ)" described the tree **before** 1A. Two packages now exist — `internal/paths` and `internal/durable` — and **no pre-existing Go file was modified**. The census of 808 is unchanged; what changed is the working tree, not the census | The recovery point in §12 |
| 6 | **1A is complete** as `TASK-01` + `TASK-02`; the next task is `TASK-03` (Phase 1B), approved 2026-11-22 as a preparation boundary only | §14's "single next owner decision: approve or amend `TASK-01`" |

### Task status, superseding the table above

| Task | State | Evidence | Next |
|---|---|---|---|
| `TASK-01` PathService | **Closed** — accepted after audit | `TASK-01-AUDIT-2026-11-22.md` §8 | none |
| `TASK-02` DurableStore | **Closed** — accepted for its bounded objective | `TASK-02-DURABLESTORE-2026-11-22.md` §8; `TASK-02-AUDIT-2026-11-22.md` | adoption — **separate task, unauthorized** |
| `TASK-03` Conversation authority | **Prepared, implementation unauthorized** | `TASK-03-CONVERSATION-AUTHORITY-PREPARATION-2026-11-22.md` | resolve D-1…D-4 into a decision record; owner review |

Census distinction: historical census **808**; current corpus **814**; placeholders `HEKI-D-0809`–`0814` are not
ledger nodes and have no recorded hashes. **814 is the corpus, not the census.**
