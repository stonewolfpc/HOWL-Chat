# C1 — P17 Evidence-Semantics Owner Decision Register (OD-1–OD-9)

**Date:** 2026-11-09
**Status:** OPEN DECISION REGISTER. **NINE DECISIONS ACCEPTED (OD-1, OD-2, OD-3, OD-4, OD-5, OD-6, OD-7, OD-8, OD-9); 0 DEFERRED; 0 UNRESOLVED.** P17 NOT APPROVED. NO REV4, P18, CONTRACT AMENDMENT, SOURCE CHANGE, TEST, BUILD, OR C1 IMPLEMENTATION.
**Inputs:** `ASSESSMENT-howl-tts-stage3-c1-p17-remaining-evidence-semantics-2026-11-09.md` (`3c4ce2be…`); P17 `01d887c1…`; corrected r3 response `60ca1cb5…`; lifecycle decision `aa136bc9…`; contract rev3 `ed0da1a7…`.
**Purpose:** record each unresolved question, its source evidence, alternatives, consequences, dependencies, a non-binding proposed default where justified, and the owner's disposition. This is a decision register, not an instruction to decide all nine now.

No frozen artifact is edited. No source, test, build, contract, or frontend file is changed. This register does not authorize any public API or provider-state change.

## 1. Reading rules

1. **The register records questions, not answers.** Every entry starts `UNRESOLVED`. The owner may mark an entry accepted, rejected, or deferred later.
2. **A proposed default is not a decision.** It is included only where the evidence supports a defensible starting point; it has no authority until the owner accepts it.
3. **Internal evidence decisions do not become public API decisions.** Any public projection or provider-state change requires a separately approved contract amendment.
4. **The thirteen-value provider-state set remains frozen.** No entry below adds, removes, or renames a public state.
5. **Owner acceptance is not independent verification.** A direction can be accepted without establishing that its implementation is correct.
6. **Identifiers are stable.** The assessment's OD-1–OD-9 numbering is retained so downstream references do not drift.

### 1.1 Identifier reconciliation

The reviewer's priority note names "OD-1 — Evidence authority", "OD-2 — Terminal authority", and "OD-3 — Meaning of `crashed`". Those labels do not match the assessment's numbering. To avoid silently renumbering the register:

| Reviewer's priority topic | Assessment/register ID | Register title |
|---|---|---|
| Evidence authority | **OD-1** | Evidence-record authority |
| Terminal authority | **OD-3** | Terminal publication authority and immutability |
| Meaning of `crashed` | **OD-5** | Meaning and scope of `crashed` |
| (not in the priority list) | **OD-2** | Instance identity |
| (not in the priority list) | **OD-4** | Public projection for `terminated_by_action` |
| (not in the priority list) | **OD-6** | Handshake compatibility taxonomy |
| (not in the priority list) | **OD-7** | Quarantine status |
| (not in the priority list) | **OD-8** | Causality standard |
| (not in the priority list) | **OD-9** | P17 approval gate |

The three priority topics are therefore **OD-1, OD-3, and OD-5** in this register. OD-2 remains instance identity. No identifier is reassigned.

## 2. Register summary

| ID | Question | Status | Proposed default (non-binding) | Downstream |
|---|---|---|---|---|
| OD-1 | Who owns the authoritative evidence record, and is it diagnostic-only or projection-supporting? | **ACCEPTED** | Accepted: projection-supporting as well as diagnostic; supervisor is the designated authority for process-lifecycle evidence and provider-state publication; client/session observations must also be represented and retained; exact sequence/synchronization/retention/conflict mechanics and contract wording remain future work | R09, R10, U1–U7 |
| OD-2 | What identity does evidence bind to: session, process, supervisor generation, or provider load? | **ACCEPTED** | Accepted: composite client session + supervisor generation + OS PID; generation is the process-lifetime key | R09, U1–U7 |
| OD-3 | Who may publish a terminal outcome, and is it immutable? | **ACCEPTED** | Accepted: immutable first terminal decision; publisher, synchronization, and late-evidence contract unresolved | R09 (terminal principle), R10, R11, U2, U3, U7 |
| OD-4 | What is the public projection, if any, for a nonzero exit after a known action? | **ACCEPTED** | Accepted: confirmed exit without orderly shutdown → `unavailable` unless a prior justified failure state applies; orderly shutdown → `absent`; active operational states cannot survive exit; `terminated_by_action` internal; no new state; action causation under OD-8 | U2, U4, U5 |
| OD-5 | Does `crashed` describe process exit, fatal helper failure, or a narrower condition? | **ACCEPTED** | Accepted: observed process exit where orderly shutdown is not established; no physical mechanism, cause, intent, or action attribution asserted; `ProcessAlive` independence is an invariant, but the five-state implementation mapping is not accepted | R13, U4, U5, U6 |
| OD-6 | Is a missing mandatory capability `incompatible` or `unavailable`? | **ACCEPTED** | Accepted: only a structurally valid, semantically identified failure of a mandatory compatibility requirement justifies `incompatible`; structural ambiguity, malformed content, and unestablished compatibility meaning remain `unavailable`; helper/client instances share the public state while retaining internal evidence; scope/extension-version categories remain open | R12, U8, U3 |
| OD-7 | Is `quarantined` internal-only, and does P17 retain a public claim for it? | **ACCEPTED** | Accepted: internal-only, not a provider state or public projection; an explicit recorded unresolved ownership/reaping condition for the exact process generation, not a timeout or `ProcessAlive` alone; later exit is not reaping; P17 claim to be corrected later; contract rev3 A3 unadopted | R11, U2, U6, U7 |
| OD-8 | What evidence establishes that an action or fatal observation caused an exit? | **ACCEPTED** | Accepted: six-dimension evidence standard; A-E required for causal attribution; A-O association only; A-U unresolved; C-unresolved/C-probable/C-established | R10, U1–U5; U5 remains OPEN |
| OD-9 | May P17 be approved with U1–U8 open? | **ACCEPTED (Option A)** | P17 is not approved while U1–U8 remain OPEN without determined outcomes or explicitly defined unknown-case rules; OD-1 is now accepted; this is a gate, not a permanent rejection | All |

**Aggregate counts (2026-11-09):** **9 ACCEPTED** (OD-1, OD-2, OD-3, OD-4, OD-5, OD-6, OD-7, OD-8, OD-9); **0 DEFERRED**; **0 UNRESOLVED**. Total: 9.

### 2.1 Review guidance received (historical; dispositions noted)

The following guidance was received on 2026-11-09. It is recorded for provenance. OD-1 through OD-9 were subsequently accepted; OD-1 and OD-4 are recorded above, and OD-5, OD-6, OD-7, OD-8, and OD-9 are recorded below. The remaining rows are guidance only, and no dependent entry is settled by this table.

| Topic | Guidance recorded | Register effect |
|---|---|---|
| OD-2 | Use a supervisor-generated instance token or generation as the authoritative process-lifetime identity. Include session identity to distinguish application lifetimes. Retain PID as corroborating diagnostic metadata, not as the identity itself. Bind observations, termination actions, exit facts, and evidence records to the same instance identity. Never let a delayed observation from an old instance modify its replacement. | This recommendation was accepted on 2026-11-09 as the OD-2 decision. Other rows are guidance unless a later disposition is noted. |
| OD-3 | Accept the immutability principle for a first terminal decision: one terminal decision per process instance, one linearization point, no later overwrite, later observations still retained against the same instance, no silent disappearance of late evidence, and no stale operation modifying a replacement instance. The exact publisher, synchronization mechanism, and late-evidence handling contract remain for later design and independent review. | Accepted on 2026-11-09 as the OD-3 principle; explicitly excluded questions remain unresolved. This does not decide OD-1's evidence-record ownership. |
| OD-3 | Distinguish an immutable terminal decision from the continuing evidence record. Later observations must remain recordable after terminal publication; whether they may change the projection is a separate decision. | Historical guidance: recordability is separate from projection mutability. OD-3 accepts the immutability principle; later-evidence effects on projections remain unresolved. |
| OD-5 | Defining `crashed` narrowly is sensible, but do not automatically substitute `degraded`, `unavailable`, or `outcome_unknown` without establishing the exact projection rules. | The proposed default is narrowed to the semantic restriction; replacement projections require a separate rule. Subsequently accepted on 2026-11-09; see the OD-5 entry and qualification. |
| OD-9 | OD-9 is a governance gate, not permission to force every unresolved U case into a deterministic outcome. An explicitly justified `outcome_unknown` rule can be a valid resolution. | Owner accepted OD-9 Option A on 2026-11-09; P17 remains unapproved. |
| Process | Work through OD-1 to OD-9 individually, beginning with evidence authority. Record the owner's actual disposition and rationale. Keep deferred decisions unresolved. Do not let one accepted decision implicitly settle a dependent one. | Dispositions are recorded one at a time; dependency links remain controlling. |

The OD-1, OD-2, OD-3, OD-4, OD-5, OD-6, OD-7, OD-8, and OD-9 decisions were accepted on 2026-11-09. There are no deferred or unresolved decisions. The OD-3 exclusions remain open.

## 3. Decision entries

Each entry uses the same fields. "Source evidence" cites the assessment's exact locations; it is not re-derived here.

### OD-1 — Evidence-record authority

| Field | Entry |
|---|---|
| **Question** | Who owns the authoritative evidence record; which components may append; and is the record diagnostic-only or allowed to support provider-state projection? |
| **Source evidence** | No evidence log or record type exists. `Failure` has only `Code`, `Message`, `Retryable`, `Field`, `Detail` (`errors.go:15–21`). `failAll` stores only the first read failure (`client.go:373–376`); `SetState` overwrites `lastFailure` (`supervisor.go:233–242`); a late terminal is dropped without a record (`client.go:339–342`); `Terminate` records no intent or result (`supervisor.go:246–288`). Assessment §2.2, §4 SRC-01–SRC-09, SRC-15–SRC-19. |
| **Alternatives** | **A. Diagnostic-only.** The log is for post-hoc diagnosis; no projection may cite it as authority. **B. Supervisor-owned, projection-supporting.** A single supervisor record is authoritative for terminal evidence and projections. **C. Client-owned, projection-supporting.** The client session record is authoritative. **D. No log.** Keep only the current status snapshot and per-failure fields. |
| **Consequences** | A narrows P17 §7 substantially and may make R09 acceptable as a diagnostic model, but it cannot resolve U1–U7. B requires a record schema, retention rule, and supervisor authority; it is the largest change and may imply contract text. C duplicates authority across client and supervisor and does not cover process exit, which the supervisor observes. D leaves P17 without an evidence model and keeps every race unresolved. |
| **Dependencies** | OD-2 (identity), OD-8 (causality), R09. Any public projection also depends on a contract amendment. |
| **Proposed default (non-binding)** | **None.** Choosing diagnostic-only is a material scope reduction that the owner must intend, not a safe default. Choosing projection-supporting presupposes schema and authority the source does not have. The owner should decide the purpose first. |
| **Owner disposition** | **ACCEPTED — projection-supporting.** Owner accepted OD-1 on 2026-11-09. The evidence record is projection-supporting as well as diagnostic, but only within explicitly justified evidence and publication rules. One accountable authority owns evidence integrity, conflict resolution, and provider-state publication; the supervisor is the designated authority for process-lifecycle evidence and provider-state publication. Client/session observations must also be represented and retained; supervisor ownership does not authorize silently discarding them. |
| **Downstream effects** | R09/R10 remain open; acceptance supplies the authority model, not the record schema, allocator, retention, or conflict implementation. U1–U7 remain OPEN. P17 remains unapproved; no contract amendment or implementation follows. |

**Accepted governance decision (2026-11-09).** The owner accepts OD-1 as **projection-supporting**. The historical question, alternatives A–D, source evidence, consequences, dependencies, and proposed default above are retained for provenance; they are not rewritten into a claim that the source already implements this model.

**Purpose.** The evidence record is projection-supporting as well as diagnostic. It may support provider-state projection, but only within explicitly justified evidence and publication rules.

**Ownership and authority.**

- One accountable authority owns evidence integrity, conflict resolution, and provider-state publication.
- The supervisor is the designated authority for process-lifecycle evidence and provider-state publication.
- Client/session observations must also be represented and retained; supervisor ownership does not authorize silently discarding them.

**Mandatory governance constraints.**

1. Every relevant record is bound to the process generation under OD-2.
2. Event order and observation/processing order remain distinct. Delayed delivery cannot establish causal ordering.
3. OD-3's immutable first-terminal-decision principle applies. Later evidence must still be retained.
4. OD-8 governs causal attribution. Temporal correlation alone is insufficient.
5. Missing evidence must not automatically mean that an event did not occur.
6. Only sufficiently justified evidence may support public provider-state projection.
7. Internal evidence does not create additional public provider states.

**Preserved boundaries.** This accepts the evidence-record authority model; it does not approve P17. OD-9 remains accepted as Option A and P17 remains unapproved. Contract rev3 remains frozen and unamended. U1–U8 remain OPEN. OD-2 through OD-9 remain accepted and unchanged. Exact sequence allocation, synchronization, retention mechanics, conflict implementation, and contract wording remain future work. No source change, test, build, benchmark, runtime experiment, or C1 implementation is authorized.

**Future work.** The authority model is fixed; the record schema, generation binding, event ordering, late-evidence reconciliation, conflict resolution, retention, and publication mechanics remain to be specified and reviewed. This decision does not close R09, R10, or any U case.

### OD-2 — Instance identity

| Field | Entry |
|---|---|
| **Question** | What is the "instance" to which observations, actions, and exit facts are bound: the client session, the OS process, a supervisor generation, or a provider load? |
| **Source evidence** | `session` is a random client string (`protocol.go:259–265`) rotated on restart (`client.go:206–223`). The process interface exposes a PID, but `Failure` carries none. The supervisor stores a `proc` pointer and one `waitErr` (`supervisor.go:166–195`), with no generation counter. `providerID`/`engine` identify a selected provider, not a process lifetime. Assessment §2.2, §4 SRC-02, SRC-09, SRC-22. |
| **Alternatives** | **A. Client session.** Simple, already generated; does not survive a restart and does not identify the process. **B. OS PID.** Identifies the process, but PIDs can be reused and the value is not retained in failures. **C. Supervisor generation.** A monotonic counter for each process lifetime, independent of PID reuse; requires a new field. **D. Composite.** Session + generation + PID, with generation as the process-lifetime key. **E. Provider load.** Binds to a load operation, not a process; insufficient for exits. |
| **Consequences** | A cannot distinguish two processes in one client lifetime. B risks attributing a reused PID. C solves process lifetime but needs an allocator. D is the most precise and makes the missing generation explicit. E conflates provider state with process ownership. |
| **Dependencies** | OD-1; R09's sequence allocator. |
| **Proposed default (now accepted)** | Composite **session + supervisor generation + PID**, with generation as the process-lifetime key. Accepted on 2026-11-09 as the OD-2 decision. This is a conceptual identity model, not an implementation instruction. |
| **Owner disposition** | **ACCEPTED** — owner accepted the generation-based model on 2026-11-09. Rationale: generation is the authoritative process-lifetime identity; session identity provides the outer application scope; PID is corroborating metadata only; observations, actions, exit facts, and evidence bind to the same instance; a delayed observation from an old instance must not modify its replacement. Acceptance is conceptual and does not authorize implementation or settle the then-unresolved OD-4–OD-9; OD-3 was accepted separately, and OD-8 was accepted later on 2026-11-09. |
| **Downstream effects** | R09 now has an accepted identity requirement, but its record schema, allocator, and retention rules remain unresolved. U1–U7 remain OPEN. No dependent decision is settled. |

**Accepted clarification (2026-11-09):** Generation uniqueness must survive supervisor-object replacement within the same session. A replacement supervisor must not restart its counter and reuse an earlier generation; an incarnation token or equivalent uniqueness mechanism is part of the accepted requirement. This clarification records a necessary condition of OD-2; it does not reopen the decision or authorize implementation.

### OD-3 — Terminal publication authority and immutability

| Field | Entry |
|---|---|
| **Question** | Who may publish a terminal provider outcome, and once published, may a later action or exit revise it? |
| **Source evidence** | `SetState` writes any state and overwrites `lastFailure` (`supervisor.go:233–242`). `Terminate` unconditionally sets `stopping` when a process exists (`:253–259`). `observeExit` writes `absent` or `crashed` (`:176–190`). `Start` publishes `incompatible`/`unavailable` after `Terminate` returns, including after a timeout (`client.go:248–254`). `Shutdown` clears `started`/`transport` after a timeout (`:668–674`). No terminal latch exists. Assessment §2.4, §2.6, §3.3, §3.7. |
| **Alternatives** | **A. Last-write-wins (current).** Any later writer may replace the terminal. **B. First-terminal latch.** The first terminal for the instance is immutable; later evidence is appended but does not republish. **C. Supervisor-authorized terminal with revision.** Only the supervisor publishes; revision is allowed only on specified new evidence. **D. No terminal provider state.** Terminal evidence remains internal/diagnostic; the frozen provider state is not used for it. |
| **Consequences** | A leaves U3/U7 unresolved and permits a crash to erase an incompatibility. B gives a defensible answer but requires an instance identity and a definition of "terminal." C requires a revision rule and a distinction between correction and overwrite. D avoids public mutation but may conflict with the existing `incompatible`/`crashed`/`absent` states and the decision's allowance that an already-justified state stays. |
| **Dependencies** | OD-1, OD-2, OD-5; U3, U7. Any public change requires a contract amendment. |
| **Proposed default (historical, non-binding)** | **First terminal is immutable for the process instance; later evidence is appended and must not republish.** Rationale: it stops the observed overwrite races without inventing a new public state. It still requires the owner to define which states count as terminal and whether a timeout can supersede one. |
| **Owner disposition** | **ACCEPTED** — owner accepted OD-3 on 2026-11-09. The accepted principle and exclusions are recorded below. |
| **Downstream effects** | No dependent finding or scenario is closed. R09, R10, R11, U1–U8 retain their prior dispositions; see the acceptance block below. |

**Accepted principle (2026-11-09).** The owner accepts:

1. Each process instance may acquire one immutable first terminal decision.
2. That decision is established at one well-defined linearization point.
3. Later observations cannot overwrite or replace the established terminal decision.
4. Later observations remain retainable evidence associated with the same process instance.
5. Late evidence must not silently disappear merely because terminal publication has occurred.
6. Evidence or operations associated with an old instance must never modify a replacement instance.

**Explicitly excluded from this acceptance — remain unresolved:**

- Who owns the authoritative evidence record.
- Which component is authorized to publish the terminal decision.
- The exact synchronization mechanism and linearization implementation.
- How late observations are accepted, retained, buffered, or reconciled.
- How terminal publication interacts with evidence-record lifecycle and shutdown.
- Whether and how later evidence may affect diagnostics or non-terminal projections.

**Downstream effect of acceptance:**

- **R09:** terminal immutability is an accepted architectural requirement; evidence schema, allocator, retention, and publisher authority remain unresolved.
- **R10:** projection precedence and causal inference remain unresolved.
- **R11–R13:** retain their existing dispositions.
- **U1–U8:** all remain OPEN; OD-3 does not establish their outcomes.
- **OD-1:** was accepted subsequently on 2026-11-09 as projection-supporting; this does not alter OD-3's exclusions.
- **OD-4–OD-9:** remained `UNRESOLVED` at the time of OD-3 acceptance. OD-4, OD-5, OD-6, OD-7, and OD-8 were accepted subsequently; OD-9 was accepted subsequently as Option A, and P17 remains unapproved.

Acceptance of an architectural principle is not proof that the existing implementation satisfies it.

### OD-4 — Public projection for `terminated_by_action`

| Field | Entry |
|---|---|
| **Question** | When a known deliberate action explains a nonzero exit, what public provider state, if any, should be published? |
| **Source evidence** | `observeExit` maps a nonzero exit observed while `stopping` to `processState = exited` and `providerState = absent` (`supervisor.go:176–182`). P17 defines `terminated_by_action` as internal and leaves its public projection undecided (`P17 §7`). `Terminate` records no intent or result (`supervisor.go:246–288`); results are ignored at `client.go:248`, `:654`, `:661`, `:664`, `:861`. Assessment §2.3, §3.5. |
| **Alternatives** | **A. Publish `absent` (current).** Treat the action as explaining the exit. **B. Keep `terminated_by_action` internal; publish no new provider state.** Preserve the last operationally justified state. **C. Publish `outcome_unknown` as a public state.** Requires a contract change; `unknown` is not in the frozen set. **D. Publish a new public state.** Requires a separately approved amendment. |
| **Consequences** | A asserts a clean outcome the evidence may not establish. B is conservative but leaves consumers with a stale state. C conflicts with the frozen thirteen-value set unless amended. D is the largest contract change and is not justified by the current evidence. |
| **Dependencies** | OD-3 (terminal authority), OD-5 (`crashed` meaning), OD-8 (causality); U5, U2, U4. |
| **Proposed default (non-binding)** | **B. Keep `terminated_by_action` internal and publish no new provider state; do not use `absent` as a silent substitute.** Rationale: it avoids claiming a clean outcome while respecting the frozen set. The owner must decide whether retaining the prior state is acceptable. |
| **Owner disposition** | **ACCEPTED** — owner accepted OD-4 on 2026-11-09 with the public-state policy and qualifications recorded below. Acceptance settles the governance policy only and does not close U2/U4/U5, define the exact contract wording, or authorize an implementation. |
| **Downstream effects** | U2/U4/U5 remain OPEN; U7 remains OPEN; historical failure evidence must remain available; P17 and contract rev3 remain unchanged. |

**Accepted governance decision (2026-11-09).** The owner accepts OD-4 and the following public-state policy. The historical question, alternatives A–D, source evidence, consequences, dependencies, and historical proposed default above are retained for provenance; they are not rewritten.

**Public-state policy.**

| Established condition | Public projection |
|---|---|
| Confirmed process exit, orderly shutdown not established | `unavailable`, unless a more specific prior failure classification remains independently justified |
| Independently established orderly shutdown | `absent` |
| Confirmed exit with an independently justified prior failure state | Preserve that failure classification where applicable |
| Process exit not yet established | Do not pretend the process is absent; preserve the last justified current state |
| Deliberate-action causation | Internal attribution under OD-8; no public provider state |

**Additional constraints.**

- `terminated_by_action` remains internal and has no public projection.
- No new public provider state is added; the thirteen-value set remains frozen.
- Active operational states cannot survive a confirmed process exit. In particular, `ready` must not be retained as a current claim after the process is gone.
- A nonzero exit alone proves neither orderly shutdown nor deliberate-action causation.
- Historical evidence remains available internally, including a prior failure condition and the action classification.
- `crashed` retains OD-5's accepted observed-exit/no-orderly-shutdown meaning.
- `unavailable` is the accepted fallback policy, not permission to discard a more specific, independently justified failure classification.
- Exact contract wording and implementation remain separate work.

**Open boundaries.** U2, U4, U5, U7, and U3 remain OPEN; U1–U8 remain OPEN. OD-1 was accepted subsequently as projection-supporting; OD-9 was accepted subsequently as Option A. P17 remains unapproved; contract rev3 remains unchanged. No source change, test, build, runtime verification, or implementation is authorized.

### OD-5 — Meaning and scope of `crashed`

| Field | Entry |
|---|---|
| **Question** | Does `crashed` describe an observed process termination, a fatal helper failure that may leave the process alive, or a narrower provider condition? |
| **Source evidence** | `Synthesize` publishes `TTSStateCrashed` on a write failure before `cleanupRejected` (`client.go:755–758`); the process may still be alive. `observeExit` publishes `crashed` on an unexpected exit (`supervisor.go:183–190`). `Status` derives `ProcessAlive` from both dimensions (`:324`), and `processAlive` returns false for `crashed`/`unavailable` (`:351–356`), so a live process can be reported as not alive. Assessment §2.6, §3.6, §4 SRC-11, SRC-15, SRC-17, SRC-20. |
| **Alternatives** | **A. Process-outcome only.** `crashed` requires an observed exit (or equivalent process termination). A live write failure must use another state. **B. Provider-condition failure.** `crashed` may describe a fatal helper condition even if the process is alive. **C. Split.** `crashed` = observed process exit; live transport failure = `degraded`/`unavailable` (or internal `outcome_unknown`). **D. Remove `crashed` from this path only.** Keep the state but declare the `client.go:756` write non-conforming. |
| **Consequences** | A makes `crashed` precise but requires changing the write-failure projection. B preserves current behavior but keeps the process/condition conflation. C resolves the conflation with the least new vocabulary. D is a narrow correction but leaves the broader semantics undefined. |
| **Dependencies** | OD-3, OD-4, OD-8; R13, U4, U5, U6. |
| **Proposed default (non-binding)** | **C. Restrict `crashed` to an observed process outcome; use `degraded` for a live transport failure that leaves the provider usable, `unavailable` only when justified, and `outcome_unknown` internally when neither is established.** Rationale: it separates process outcome from provider condition without adding a public state. The choice among `degraded`/`unavailable` still needs the owner. |
| **Owner disposition** | **ACCEPTED** — owner accepted OD-5 on 2026-11-09. The accepted meaning and the `ProcessAlive` qualification are recorded below. Acceptance settles the semantic question only; it does not authorize implementation or close R13, U4, U5, or U6. |
| **Downstream effects** | R13 is not closed; U4/U6 remain OPEN; U5's public projection remains OD-4; P17 and contract rev3 remain unchanged. |

**Accepted governance decision (2026-11-09).** The owner accepts OD-5 and the following meaning. The historical question, alternatives A–D, source evidence, consequences, dependencies, and historical proposed default above are retained for provenance; they are not rewritten.

**Accepted meaning.** `crashed` represents an observed process exit for which an orderly shutdown has not been established. It does not assert a physical crash mechanism, the cause of exit, intent, or deliberate-action attribution.

**Four-layer separation.**

1. Process exit is an observed lifecycle fact.
2. Orderly shutdown is a separately established outcome.
3. Deliberate-action causation requires OD-8's accepted A-E standard.
4. Provider/runtime failure does not establish process exit.

**Operational boundaries.** A live process must not be represented as `crashed`; a provider/runtime condition, a fatal transport/protocol observation, a stop request, `stopping`, or a lifecycle closure is not an observed process exit. `stopping`, a stop request, temporal ordering of observations, a successful termination API result, and a nonzero exit code do not independently establish that a deliberate action caused the exit. The source's `stopping` branch is an intent-based classification, not proof of orderly shutdown. No public provider state is added, removed, or renamed.

**`ProcessAlive` qualification.** `ProcessAlive` independence from provider state is accepted as a design invariant. The proposed five-state mapping is **not** automatically accepted as an implementation contract. In particular, `starting` must not be treated as proof of actual OS-process liveness merely because it is a lifecycle label. The authoritative meaning is actual process liveness, independently of provider state; the implementation must establish that from process ownership/handle evidence. The exact liveness predicate is not approved or implemented.

**What remains open.** U6, U4, and U5 remain **OPEN**. OD-4 was accepted subsequently. R13 is not closed. The current `client.go:756` write and the current `Status` liveness derivation are not declared conforming. P17 and contract rev3 are unchanged. No implementation authorization is granted.

### OD-6 — Handshake compatibility taxonomy

| Field | Entry |
|---|---|
| **Question** | Is a missing mandatory capability an `incompatible` condition or a `protocol_violation`/`unavailable` condition, and how should helper-declared and client-detected failures be classified? |
| **Source evidence** | `validateHello` has nine failure exits, all `protocol_violation` (`protocol.go:502–545`); missing mandatory capabilities is one (`:527–531`). `Start` maps only `protocol_incompatible`/`protocol_minor_unsupported` to `incompatible`; all other handshake failures become `unavailable` (`client.go:250–254`). `failureFromWire` preserves a helper-declared code (`client.go:434–435`, `:471–481`). Assessment §2.5, §3.8, §4 SRC-03, SRC-21. |
| **Alternatives** | **A. Capability mismatch = `incompatible`; malformed protocol = `unavailable`.** **B. All handshake validation failures = `unavailable` (current for client-detected).** **C. All capability/version mismatches = `incompatible`, including missing mandatory capabilities.** **D. Keep a distinct internal class and map to one public state only after a contract decision.** |
| **Consequences** | A gives users a meaningful compatibility signal but changes the current missing-capability projection. B preserves current behavior but treats a capability mismatch as generic unavailability. C may overstate incompatibility for a helper that is merely incomplete. D preserves evidence but leaves the public mapping unresolved. |
| **Dependencies** | OD-3, OD-5; R12, U8, U3. Any public change requires a contract amendment. |
| **Proposed default (non-binding)** | **A. Capability-set mismatch (including missing mandatory capabilities) → `incompatible`; malformed protocol/schema → `unavailable`.** Rationale: it aligns the projection with the compatibility concept while keeping malformed protocol distinct. The owner must confirm that a missing capability is genuinely an incompatibility, not a protocol violation. |
| **Owner disposition** | **ACCEPTED** — owner accepted OD-6 on 2026-11-09. The accepted principle and limitations are recorded below; acceptance settles the governance classification only and does not authorize implementation or close R12, U8, or U3. |
| **Downstream effects** | R12 is not closed; U8 remains OPEN; U3 remains OPEN; the exact contract definitions and internal code remain open. P17 and contract rev3 remain unchanged. |

**Accepted governance decision (2026-11-09).** The owner accepts OD-6 and the following governance principle. The historical question, alternatives A–D, source evidence, consequences, dependencies, and historical proposed default above are retained for provenance; they are not rewritten.

**Accepted principle.** Only a structurally valid, semantically identified failure to satisfy a mandatory compatibility requirement justifies `incompatible`. Structural ambiguity, malformed content, and failures whose compatibility meaning cannot be established remain `unavailable`.

**Public projection and evidence.** Helper-declared and client-detected instances of the same established mandatory compatibility failure share the public `incompatible` projection while retaining distinct internal diagnostic evidence. No new public provider state is added.

**Categories within the accepted principle.**

- Missing mandatory service capability → `incompatible`.
- Unknown required core capability → `incompatible`.
- Known required capability at an unsupported required version → `incompatible`.
- Duplicate identities → `unavailable` (structural ambiguity).
- Invalid semantic descriptors → `unavailable`.
- Malformed helper version or limits mismatch → `unavailable`.

**Categories outside the accepted rule.** Misplaced descriptors, unsupported extension versions, unknown optional/provider feature requiredness, and any classification not listed above remain unclassified. The current source does not validate descriptor placement or a supported-extension-version registry. This acceptance does not turn those gaps into an undocumented scope rule or a contract definition. A category outside the rule is not assigned to `incompatible` by default.

**Contract uncertainty.** The exact definitions of a mandatory capability, a valid descriptor, a supported required version, and descriptor placement remain open. This acceptance does not adopt A2-D3 or A2-D4, does not amend P17 or contract rev3, and does not add a public state.

**Implementation uncertainty.** The internal failure code for a client-detected compatibility failure remains open. No source change, implementation, test, build, or runtime verification is authorized.

**U8 and U3 remain OPEN.** R12 is not closed. P17 and contract rev3 remain unchanged. OD-1 was accepted subsequently as projection-supporting; OD-9 was accepted subsequently as Option A; P17 remains unapproved.

### OD-7 — Quarantine status

| Field | Entry |
|---|---|
| **Question** | Is `quarantined` an internal ownership condition only, and does P17 retain any public projection claim for it? |
| **Source evidence** | Neither the five-value process-state set nor the thirteen-value provider-state set contains `quarantined` (`internal/audio/types/provider.go:8–36`). `Terminate` can return `timeout` after a second grace while the process may remain owned (`supervisor.go:277–286`). `Shutdown` clears client transport regardless (`client.go:668–674`). The lifecycle decision says `quarantined` is not a provider state and is an internal ownership condition (`DECISION…`, §1.3–1.4). P17 §10 still lists it among five projections, while §7 has no path. Assessment §2.4, §3.2, §3.7. |
| **Alternatives** | **A. Internal-only; no public projection.** Record ownership uncertainty in an internal disposition; P17 removes the public claim. **B. Public provider state.** Requires a contract amendment; conflicts with the decision's current direction. **C. Remove the concept.** Treat a terminate timeout as an ordinary action failure with no quarantine record. **D. Internal-only with a diagnostic exposure.** Keep it out of the provider enum but expose it as a non-state diagnostic under a separate, approved contract. |
| **Consequences** | A matches the existing decision and resolves the unreachable-projection defect by narrowing P17. B reopens the frozen taxonomy. C loses the ownership distinction the decision created. D requires a separate exposure decision but preserves the internal concept. |
| **Dependencies** | OD-1, OD-3; R11, U2, U6, U7. |
| **Proposed default (non-binding)** | **A. Internal-only; no public projection. P17's public projection claim is removed, and the internal disposition records ownership without becoming a provider state.** Rationale: it is consistent with the existing owner decision and requires no enum change. |
| **Owner disposition** | **ACCEPTED** — owner accepted OD-7 on 2026-11-09 with the qualification recorded below. Acceptance settles the governance classification only and does not close R11, U2, U6, or U7, define the internal transition, or authorize an implementation. |
| **Downstream effects** | R11 is narrowed but not closed; U2/U6/U7 remain OPEN; the internal record/transition, retention, and release rules remain open; P17 and contract rev3 remain unchanged. |

**Accepted governance decision (2026-11-09).** The owner accepts OD-7 and the following principle. The historical question, alternatives A–D, source evidence, consequences, dependencies, and historical proposed default above are retained for provenance; they are not rewritten.

**Internal-only status.** `quarantined` is an internal supervisor ownership condition. It is **not** a fourteenth provider state, is not added to the frozen thirteen-value set, and has no public projection.

**Meaning of quarantine.** Quarantine means that the supervisor has **explicitly recorded an unresolved ownership/reaping condition for the exact process generation**. It is not merely a timeout and not merely a claim that a process is alive.

**Three distinctions that are not interchangeable.**

1. **Termination timed out.** The bounded termination path did not observe completion. A timeout alone does not establish quarantine.
2. **The supervisor still owns the process instance.** Retained ownership is relevant evidence, but ownership alone does not establish that the instance is safely retained or that its eventual cleanup is accounted for.
3. **The supervisor has established that the instance is safely retained and its eventual cleanup is accounted for.** This is the ownership/reaping condition quarantine denotes. The exact evidence, transition, retention, and release rules remain open.

**Identity and release.** OD-2 identity applies. A quarantine record belongs to the exact process generation it concerns and must never affect a replacement process. A later exit does not automatically establish successful reaping; ownership release requires its own evidence.

**P17 and contract rev3.** P17's contradictory public projection claim must eventually be corrected, but acceptance is not permission to revise P17 now. P17 §7's omission is consistent with internal-only status; P17 §10's claim is not. Contract rev3 §A3 remains an unadopted proposal, including its state machine, retained `Wait` ownership, and "reported distinctly" language.

**Open boundaries.** R11 is narrowed but not closed. U2, U6, and U7 remain OPEN; U1–U8 remain OPEN. The exact quarantine transition, retention period, release/reaping ownership, and any diagnostic exposure remain open. OD-1 was accepted subsequently as projection-supporting; OD-4 was accepted subsequently; OD-9 was accepted subsequently as Option A; P17 remains unapproved. No public or diagnostic exposure is authorized; any exposure requires a separate approved contract amendment.

**No implementation.** No source change, test, build, runtime verification, or P17/contract edit is authorized by this acceptance.

### OD-8 — Causality standard

| Field | Entry |
|---|---|
| **Question** | What evidence is sufficient to establish that a termination action or a fatal observation caused an exit, rather than merely preceding it? |
| **Source evidence** | `protocolViolation` calls `failAll` then `Kill` (`client.go:364–367`); `Terminate` sets `stopping` and waits (`supervisor.go:246–288`); `observeExit` decides only from `providerState == stopping` (`:176–190`). There are no timestamps, sequence numbers, action IDs, or exit-cause records. `waitErr` is an error, not an exit code plus cause. Assessment §2.3, §3.1–§3.5. |
| **Alternatives** | **A. No causal claim without action record + exit fact + ordering; otherwise `outcome_unknown`.** **B. Treat action-before-exit as sufficient.** **C. Require action completion and exit status, but not a global timestamp.** **D. Require independent external evidence (for example, OS exit status and a helper-side record).** |
| **Consequences** | A is the most conservative and keeps unknown cases honest. B would close cases by assumption and repeat the defect R10 identifies. C is a middle ground but still needs a definition of ordering. D is strongest but depends on evidence the source does not currently retain. |
| **Dependencies** | OD-1, OD-2, OD-3; R10, U1–U5. |
| **Proposed default (historical, non-binding)** | **A. No causal claim without an action record, an exit fact, and a recorded ordering; otherwise `outcome_unknown`.** Rationale: it matches the assessment's conclusion that correlation is not causation and avoids inventing a winner. The accepted standard below retains this as a floor but adds event-order, mechanism-operation, and competing-explanation requirements. |
| **Owner disposition** | **ACCEPTED** — owner accepted the OD-8 evidence standard on 2026-11-09. The accepted standard, classifications, lifecycle boundary, and U5 reservation are recorded below. |
| **Downstream effects** | Before acceptance, U1–U5 could not be closed without a standard, and R10 remained a selection policy. After acceptance, U1–U5 remain OPEN; U5 remains unresolved on causation; R10 is still not a resolution. OD-4a and the OD-3 exclusions remain unresolved; OD-4, OD-5, OD-6, and OD-7 were accepted subsequently; OD-4b's public projection is governed by OD-4. |

**Accepted governance decision (2026-11-09).** The owner accepts OD-8 and the following evidence standard. The historical question, alternatives A–D, source evidence, consequences, dependencies, and historical proposed default above are retained for provenance; they are not rewritten. The accepted standard is stricter than historical alternative A: it retains A's no-causal-claim floor and adds event-order, mechanism-operation, and competing-explanation requirements. Basis: the owner governance decision of 2026-11-09; the decision-ready proposal is `f2af1800…`.

**Rationale.** Temporal association alone must never be represented as causal attribution. A conservative standard preserves the A/B/C distinction, binds evidence to OD-2's instance identity, requires reliable event order and mechanism operation, and leaves genuinely ambiguous cases unresolved. It does not force a causal conclusion because a lifecycle decision is required.

**A. Process-instance identity.** The action, action outcome, exit fact, and supporting evidence must be bound to the same authoritative process instance. Use OD-2's accepted generation-based identity principle, including uniqueness across supervisor-object replacement within a session. PID alone is insufficient. Evidence from an old process instance must not be attributed to a replacement instance.

**B. Action identity and outcome.** Identify the specific deliberate action and its target. Distinguish requested, accepted, attempted, completed, rejected, timed out, and failed or otherwise unsuccessful. Do not treat a request as proof of execution, or a successful API return as proof of causal effect.

**C. Exit fact.** Identify the actual process exit observation and preserve available exit information, including status, signal, error, observation context, and instance identity where available. Do not substitute a generic `waitErr`, provider state, or lifecycle label for the actual exit evidence.

**D. Event ordering.** Use these classifications:

- **A-E — Event order established:** reliable evidence establishes the relevant action effect occurred before the process exit.
- **A-O — Observation order only:** the action record precedes the exit observation, but actual event order is not established.
- **A-U — Unreliable or unknowable order:** no reliable event or observation ordering can be established.

A-E is required for any causal attribution, including C-probable. A-O establishes association only; it cannot support C-probable or C-established. A-U remains unresolved and cannot support causal attribution. A sequence number is authoritative only if allocated at the relevant event occurrence or acceptance point. A number assigned when delayed observations are processed establishes observation order, not event order. Wall-clock timestamps alone must not be presumed to establish event order.

**E. Causal mechanism.** A causal attribution requires evidence of:

1. A mechanism capable of producing the observed exit.
2. Evidence that the mechanism actually operated on this particular process instance.
3. Consistency between the mechanism and the observed exit facts.
4. An explicit assessment of material competing explanations.
5. An explicit statement of remaining uncertainty.

A plausible mechanism narrative is insufficient. A successful action result, a kill-consistent exit status, `providerState == stopping`, or temporal correlation is not independently sufficient to establish causation. Do not use the same state observation as circular proof of action identity, event ordering, and mechanism operation.

**F. Competing explanations and defeaters.** Evaluate material competing explanations using available evidence. Record:

- the alternative explanation;
- evidence supporting or contradicting it;
- whether it is excluded, materially weakened, still plausible, or unresolved;
- what evidence is missing.

Do not manufacture a defeater merely to satisfy a checklist. If no material competing explanation is identified, document the scope of the assessment and why none was identified. Missing evidence for a competing explanation is not evidence in favor of the preferred causal explanation. A material unresolved alternative limits the strength of the causal conclusion.

**Causal attribution classifications.** Preserve these distinctions:

- **C-unresolved:** available evidence does not support a causal conclusion.
- **C-probable:** A-E is established, the action mechanism is evidenced as operating on the instance, exit facts are consistent, and the evidence favors action-caused termination over evaluated alternatives, but material uncertainty remains.
- **C-established:** A-E is established, the mechanism's operation is supported by strong convergent evidence, exit facts are consistent, and no material competing explanation remains unresolved.

These classifications describe the strength of the evidence. They are not public provider states and do not amend the frozen thirteen-value provider-state contract. Do not assign numerical probabilities without an independently justified statistical basis.

**Lifecycle closure is separate.** A process-lifecycle record may be closed when its termination is established, even when the cause remains unknown. Lifecycle closure must never be used as evidence that a deliberate action caused the exit. The final causal classification may remain C-unresolved. Late evidence may strengthen, weaken, or resolve a causal attribution, but it must not rewrite OD-3's immutable first terminal decision. OD-3's unresolved exclusions remain unresolved, including evidence-record ownership, publisher authorization, synchronization/linearization implementation, late-evidence reconciliation, and evidence lifecycle/shutdown interaction.

**U5 remains OPEN.** OD-8 defines how U5 must eventually be evaluated; it does not establish that U5 satisfies the standard. The existing source's sequence of a nonzero exit after a deliberate kill remains insufficient to establish causation without the required same-instance evidence, action outcome, reliable event order, mechanism-operation evidence, and competing-explanation assessment. This distinction is recorded in the downstream implications.

**Downstream implications of acceptance.**

- **OD-4a:** the standard now defines what would be needed to show that an action caused an exit; no particular action/exit case is decided.
- **OD-4b:** public projection is now governed by OD-4's accepted policy; the exact action-causation evidence remains open.
- **OD-5:** accepted subsequently on 2026-11-09 with the meaning and `ProcessAlive` qualification in its entry. OD-8 does not decide whether a live process may be `crashed`; OD-5 now does, while U6, U4, and U5 remain open; OD-4 was accepted subsequently.
- **OD-6:** accepted subsequently on 2026-11-09 with the mandatory-compatibility principle in its entry; U8 and U3 remain OPEN, and the contract/implementation details remain open.
- **OD-7:** accepted subsequently on 2026-11-09 with the internal-only ownership/reaping meaning in its entry; R11 is narrowed but not closed, U2/U6/U7 remain OPEN, and P17/contract rev3 remain unchanged.
- **OD-1:** accepted subsequently as projection-supporting; OD-8 does not supply the record schema, allocator, retention, or publication mechanics. U1–U7 remain OPEN.
- **OD-3:** the immutable-terminal principle remains accepted; its publisher, synchronization, late-evidence, and lifecycle exclusions remain unresolved.
- **R10:** the causal standard now exists, but R10 is still not a demonstrated resolution of competing evidence; selection order remains a policy.
- **U1–U8:** all remain OPEN. OD-8 does not close any U case.
- **U5:** remains OPEN. The accepted standard defines how U5 must be evaluated but does not establish that the existing nonzero-exit-after-deliberate-kill sequence satisfies it; the required same-instance, action-outcome, event-order, mechanism-operation, and competing-explanation evidence remains missing.
- **P17:** remains unapproved. OD-9 was accepted subsequently as a gate (Option A); it does not approve P17.
- **Contract:** the thirteen-value provider-state set remains frozen; no public state is added or changed.

### OD-9 — P17 approval gate

| Field | Entry |
|---|---|
| **Question** | May P17 be approved while U1–U8 remain open, or must every case have a determined outcome or an explicit `outcome_unknown` rule before approval? |
| **Source evidence** | P17 is not approved; U1–U8 are open; R01–R13 retain individual dispositions (`P17 §11`; assessment §8). The review process has repeatedly refused to close U cases by assigning higher precedence to one event. |
| **Alternatives** | **A. No approval until OD-1–OD-8 are dispositioned and every U case has a determined outcome or an explicit unknown rule.** **B. Approve P17 as a partial model with open U cases.** **C. Split P17: approve a minimal, fully specified subset and defer the rest.** **D. Defer P17 entirely until rev4/P18 work is authorized.** |
| **Consequences** | A preserves the evidence standard but may delay any progress. B risks approving an ambiguous model. C could allow a bounded, useful subset but requires a scope decision. D is a full stop and may leave the lifecycle work without an evidence model. |
| **Dependencies** | OD-1–OD-8; all R and U dispositions. |
| **Proposed default (non-binding)** | **A. No approval until OD-1–OD-8 are dispositioned and each U case has either a determined outcome or an explicit `outcome_unknown` rule.** Rationale: approval should mean the model is dependable within its declared scope, not that the open cases have been forgotten. |
| **Owner disposition** | **ACCEPTED — Option A.** Owner accepted OD-9 on 2026-11-09. P17 remains unapproved while U1–U8 remain OPEN without determined outcomes or explicitly defined unknown-case rules. OD-1 was accepted subsequently; that satisfies the OD-1 condition but does not satisfy the U-case condition. This is a governance gate, not a permanent rejection. |
| **Downstream effects** | P17 remains unapproved; OD-1 is accepted; U1–U8 remain OPEN; contract rev3 remains frozen. No contract amendment or implementation follows. |

**Accepted governance decision (2026-11-09).** The owner accepts OD-9, **Option A**. The historical question, alternatives A–D, source evidence, consequences, dependencies, and historical proposed default above are retained for provenance; they are not rewritten.

**Gate.** P17 is not approved while:

- U1–U8 remain `OPEN` without determined outcomes or explicitly defined unknown-case rules.

OD-1 is now `ACCEPTED`; its acceptance removes the OD-1 condition but does not satisfy the U1–U8 condition or approve P17.

This is a governance approval gate; it is not a permanent rejection of P17. Reconsideration of P17 approval requires the outstanding dependencies to be resolved or addressed through explicit, separately justified governance decisions.

**Preserved boundaries.** Acceptance of OD-9 does not approve P17. It does not amend contract rev3, which remains frozen. It does not authorize C1 implementation, source changes, tests, builds, benchmarks, or runtime experiments. OD-1 is `ACCEPTED`; U1–U8 remain `OPEN`; OD-2 through OD-8 remain accepted and unchanged. No P17, contract, assessment, source, or implementation change is authorized by this acceptance.

**No implementation.** P17 remains an unapproved proposal. No source change, test, build, benchmark, runtime experiment, contract amendment, or C1 implementation is performed or authorized.

## 4. Publication and amendment guardrail

1. The thirteen provider states remain exactly: `absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping`, `crashed`, `incompatible`, `unavailable`, `backoff`.
2. No entry in this register authorizes a new public state, a renamed state, or a public projection.
3. An internal decision about evidence, disposition, or quarantine does not become a public API decision by implication.
4. Any public exposure of a new distinction requires a separately approved contract amendment with a defined projection.
5. A decision to keep the current behavior is still a decision; it must be recorded explicitly, not inferred from silence.
6. OD-8's `C-unresolved`, `C-probable`, and `C-established` are evidence classifications, not public provider states. They do not amend the thirteen-value set.

## 5. Relationship to independent review

Owner acceptance of a direction is not independent verification. The following remain separate:

- **Owner decision:** what guarantee the system is supposed to provide.
- **Independent review:** whether the proposed model, source behavior, and implementation actually provide it.
- **Implementation evidence:** tests, runtime observation, and platform checks, none of which this register authorizes.

An entry marked `accepted` later records the owner's direction only. It does not close R09–R13 or U1–U8, and it does not make P17 approved. OD-8's acceptance is an evidence-governance direction; it is not independent verification and does not close R10 or U5. OD-5's acceptance is likewise a semantic direction; it is not independent verification and does not close R13 or U6. OD-6's acceptance is a governance-classification direction; it is not independent verification and does not close R12, U8, or U3. OD-7's acceptance is an internal-ownership direction; it is not independent verification and does not close R11, U2, U6, or U7. OD-4's acceptance is a public-state policy direction; it is not independent verification and does not close U2, U4, or U5. OD-9's acceptance is a gate direction; it does not approve P17 or close U1–U8. OD-1's acceptance is an evidence-authority direction; it does not close R09, R10, or U1–U7.

## 6. Current boundary

- P17 remains **not approved**.
- R01–R13 retain their existing individual dispositions.
- E1–E9 remain individually represented.
- U1–U8 remain OPEN.
- The thirteen-value provider-state set remains frozen.
- No rev4, P18, contract amendment, source change, test, build, benchmark, or C1 implementation is authorized.
- This register contains **nine accepted decisions (OD-1, OD-2, OD-3, OD-4, OD-5, OD-6, OD-7, OD-8, OD-9)**; there are no deferred or unresolved decisions.
