# C1 / P17 — OD-7 `quarantined` Status Assessment (Read-Only)

**Date:** 2026-11-09

**Status:** Assessment only. **OD-7 remains UNRESOLVED.** No register update, disposition, contract amendment, P17 edit, source change, or implementation is made or recommended by this document.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `3bdc3386…`, 350 lines.
- Five decisions are accepted: OD-2, OD-3, OD-5, OD-6, OD-8. OD-1 is deferred; OD-4, OD-7, and OD-9 are unresolved.
- OD-6 acceptance is a governance classification only; R12, U8, and U3 remain open.
- P17 remains unapproved. Contract rev3 is a frozen proposal, not an accepted amendment.
- Lifecycle decision: `aa136bc9…`; `quarantined` is not a provider state and is an internal ownership condition.
- OD-1 remains deferred; OD-9 is untouched.

## 1. Exact OD-7 question and current status

The register asks:

> **Is `quarantined` an internal ownership condition only, and does P17 retain any public projection claim for it?**

The recorded status is **UNRESOLVED**. The alternatives are:

- **A. Internal-only; no public projection.** Record ownership uncertainty in an internal disposition; P17 removes the public claim.
- **B. Public provider state.** Requires a contract amendment; conflicts with the frozen thirteen-value set and the lifecycle decision.
- **C. Remove the concept.** Treat a termination timeout as an ordinary action failure with no quarantine record.
- **D. Internal-only with diagnostic exposure.** Keep it out of the provider enum but expose it as a non-state diagnostic under a separate approved contract.

The non-binding default is A. This assessment tests that default against the source, the lifecycle decision, P17, and contract rev3.

## 2. Source inventory

| Location | Source behavior | Relevance |
|---|---|---|
| `internal/audio/types/provider.go:8–16` | Process states are `absent`, `starting`, `alive`, `stopping`, `exited`. No `quarantined`. | The process-state dimension has no quarantine value. |
| `types/provider.go:20–36` | Provider states are the frozen thirteen; no `quarantined`. | The public provider-state set has no quarantine value. |
| `types/provider.go:57–61` | `TTSProviderStatus` carries process state, provider state, and a derived `ProcessAlive`. | The observable surface has no quarantine field. |
| `supervisor.go:246–288` | `Terminate` sets process and provider state to `stopping`, closes the process, waits for grace, kills on timeout, then waits a second grace. A second timeout returns `timeout` without clearing `s.proc`. | A process can remain supervisor-owned after termination reports timeout. |
| `supervisor.go:277–286` | The second timeout is the `timeout` failure path. | This is the natural quarantine-candidate path, but the timeout is not proof of liveness. |
| `supervisor.go:291–293` | `Kill` delegates to `Terminate` with the default grace. | The same ownership result applies to the prompt-kill path. |
| `supervisor.go:306–333` | `Status` copies process and provider state and computes `ProcessAlive`. | After a timeout the status can still show `stopping` and a live process; there is no quarantine report. |
| `supervisor.go:335–343` | `Supervisor.ProcessAlive()` is true when `proc != nil` and process state is not `exited`. | The supervisor tracks process ownership separately from provider state, but not as quarantine. |
| `supervisor.go:118–124` | `Start` returns the existing `s.proc` when one is installed. | A retained process is reused rather than classified; there is no admission gate tied to quarantine. |
| `supervisor.go:233–242` | `SetState` can overwrite provider state while a process is still owned. | Provider state alone cannot encode ownership uncertainty. |
| `client.go:645–674` | `Shutdown` sends shutdown, calls `Terminate`, then sets `closing = true`, `started = false`, and `transport = nil` regardless of the termination result. | After a timeout, the client can report not started while the supervisor still owns a process. |
| `client.go:206–214` | `Start` refuses a closing client with a `shutdown` failure. | Clearing the transport does not by itself release the retained process or create a quarantine record. |
| `client.go:859–864` | `cleanupRejected` checks `sup.ProcessAlive()` and calls `Terminate` when true. | The client already has a process-ownership cleanup path, but no quarantine transition. |

### 2.1 Source-verified behavior

1. **No public quarantine state exists.** The provider-state set has thirteen values, none of which is `quarantined`. The process-state set has five values, none of which is `quarantined`.
2. **A timeout can leave the process owned.** `Terminate` returns `timeout` after the second grace without clearing `s.proc`; `ProcessAlive` can remain true.
3. **A timeout does not prove the process is alive or dead.** It proves that the termination path did not observe completion within its bounds. The process may have exited without the observer completing, or may still be running. The exact ownership/liveness fact requires separate evidence.
4. **No internal quarantine record exists.** There is no supervisor field, transition, or status field that records `quarantined`, a retained `Wait` owner, or an ownership-uncertainty disposition.
5. **`Shutdown` can separate client and supervisor views.** The client clears `started` and `transport` regardless of the termination error; the supervisor may still own the process. The client then reports not started, while the supervisor's process view can remain `stopping`/alive.
6. **Provider state cannot carry the concept.** `SetState` can overwrite provider state; `ProcessAlive` is derived; neither is a quarantine value. A provider state must not be used as a substitute.

## 3. What P17 currently claims

P17 §10 challenge 5 asks the reviewer to confirm five required projections: `absent`, `crashed`, `incompatible`, `unknown`, and `quarantined` (`PROPOSAL…`, §10). P17 §7's projection table has no `quarantined` row and no path that produces it.

The claim is therefore inconsistent in two ways:

- it names `quarantined` as a required projection while §7 omits it;
- the source has no `quarantined` value or transition.

P17 is unapproved. This assessment does not edit it. The correct OD-7 direction, if A is accepted, is to remove or narrow the §10 claim in a future P17 revision: `quarantined` is not a public projection. The §7 omission is consistent with that direction, but §7 alone does not establish the internal disposition.

## 4. What the lifecycle decision already settles

The lifecycle decision (`DECISION…`, §1.3–1.4) states:

- `quarantined` is **not** a provider state and is not added to the frozen contract;
- provider state and process-lifecycle disposition are two dimensions;
- the internal lifecycle dimension is conceptually `ProcessDisposition { active, reaping, reaped, quarantined, outcome_unknown }`, with names expressly illustrative rather than an implementation instruction;
- public exposure of a new distinction requires a separately approved contract amendment with a defined projection.

This supports the **internal-only** classification in OD-7-A. It does not define:

- the exact predicate that makes a process quarantined;
- who records it, when it is cleared, or how long it is retained;
- how the retained process is reaped or released;
- whether any diagnostic exposes it.

Those remain open. OD-7 cannot treat the illustrative names as an implementation contract.

## 5. Relationship to contract rev3's quarantine proposal

Contract rev3 §A3 proposes a richer internal state machine:

- `pending → live → reaping → reaped`, with `reaping → quarantined`;
- a quarantine entry retaining `proc`, lifecycle, `waitOwner`, `lastErr`, and `attempts`;
- a bounded escalation ladder, no second kill, retained `Wait` ownership, and no admission lock;
- a statement that a quarantined process is "reported distinctly from `absent` and from `crashed`."

That is a **proposal**, not accepted source behavior or an adopted contract. Its "reported distinctly" language is ambiguous: it could mean internal diagnostic reporting or a public projection. It cannot be used to settle OD-7 in favor of a public state. It also cannot be used to claim that quarantine is implemented.

The proposal is useful as design evidence for A or D, but it does not establish:

- that the source has a quarantine record;
- that a timeout alone is sufficient to enter quarantine;
- that a quarantined process is publicly reported;
- that the proposed retained-`Wait` ownership is compatible with the current supervisor.

## 6. Alternatives assessed

| Alternative | Source fit | Public consequence | Assessment |
|---|---|---|---|
| **A. Internal-only; no public projection** | Matches the lifecycle decision and the absence of any public `quarantined` value. The source has ownership evidence but no quarantine record; an internal disposition would need to be added later. | Keeps the thirteen-value set frozen; P17's public claim is removed; no new public state. | **Recommended.** It preserves the distinction without inventing a public state or a source capability. |
| **B. Public provider state** | Conflicts with the frozen enum and the lifecycle decision. Requires a contract amendment. | Adds a fourteenth public state. | Rejected for now. It would turn an internal ownership condition into a provider condition without a defined public meaning. |
| **C. Remove the concept** | Treats a timeout as an ordinary action failure. Loses the distinction between "termination did not complete" and "process still owned." | No new state; simplest surface. | Not recommended. It discards the ownership distinction the lifecycle decision created and leaves retained processes unrepresented internally. |
| **D. Internal-only diagnostic exposure** | Preserves the internal concept; exposure requires a separate approved contract. | No provider state; may add a non-state diagnostic. | Not selected here. It is compatible with A but adds an exposure decision OD-7 does not need to make. |

### 6.1 Why A is the narrowest defensible decision

A confirms only what the lifecycle decision and the source support:

- there is no public `quarantined` provider state;
- process ownership can outlive a client's transport view;
- that ownership condition is internal, not a provider condition;
- P17 should not claim a public projection the model cannot produce.

A does **not** claim that the source already records quarantine, that a timeout proves quarantine, or that retained process ownership is safely reaped. Those remain implementation and contract questions.

## 7. Recommended OD-7 decision formulation

The following is a recommendation for owner consideration, not a recorded decision.

> **OD-7 — `quarantined` status (proposed):**
>
> 1. `quarantined` is **not** a provider state and has no public projection. The thirteen-value provider-state set remains frozen.
> 2. `quarantined` is an internal supervisor ownership condition: a process may still be owned after a termination/reaping sequence has not established release. The condition is distinct from `absent`, `crashed`, `unavailable`, and `outcome_unknown`.
> 3. A termination timeout is a **candidate**, not sufficient proof, for quarantine. The internal predicate must use process-ownership/liveness evidence; a timeout alone does not establish that the process is alive, dead, or safely retained.
> 4. The internal condition must bind to the OD-2 process instance and must not be inferred from provider state, `ProcessAlive` alone, or a later exit alone. Clearing quarantine requires evidence that reaping/release completed; a later exit does not automatically clear it.
> 5. P17's public projection claim is removed or narrowed: `quarantined` is not one of the required public projections. P17 §7's omission is consistent with this, but no P17 edit is made by OD-7.
> 6. Contract rev3 §A3 is not adopted by this decision. Its proposed state machine, retained `Wait` ownership, and "reported distinctly" language require separate review.
> 7. Any public or diagnostic exposure of quarantine requires a separately approved contract amendment with a defined projection. OD-7 selects internal-only, not diagnostic exposure.
> 8. OD-7 does not decide the internal record fields, transition predicate, reaping ownership, retention period, or recovery action.
> 9. OD-7 does not decide R11 closure, U2, U6, U7, or U3. Those remain open.

### 7.1 What it would make non-conforming

If accepted, any P17 statement that `quarantined` is a required public projection would be non-conforming and would need a future P17 revision. The current source has no public `quarantined` value, so it does not currently publish one; the source also has no internal record, so it does not currently conform to the internal-disposition part of the recommendation.

The decision would not require a new public state, a provider enum change, or a source change by itself.

### 7.2 What it would not decide

- The exact internal representation of quarantine.
- Whether quarantine is entered on a termination timeout, a failed reaping step, or another ownership condition.
- Whether the supervisor retains `Wait` ownership, and who may reap a quarantined process.
- Whether a client diagnostic may expose the internal condition without a contract amendment.
- Whether a later exit or a successful cleanup clears quarantine.
- Whether R11 is closed; it can be narrowed but not closed without the internal path and a P17 correction.
- U2, U6, U7, and U3.

## 8. Consequences

### 8.1 R11

R11 is a design omission: P17 claims a projection that §7 and the source do not support. OD-7-A resolves the public claim by confirming internal-only status. It does not close R11, because the internal record/transition and the retained-process lifecycle are still undefined.

### 8.2 U7

U7 asks about a `Terminate` timeout after the terminal projection is published. OD-7-A keeps quarantine internal and does not decide whether the terminal state is immutable or how a later ownership condition affects it. U7 remains OPEN.

### 8.3 U2

U2 concerns a shutdown timeout racing a protocol violation. OD-7-A does not choose a winner between the two failures and does not turn quarantine into a public state. U2 remains OPEN.

### 8.4 U6

OD-5's accepted `crashed` meaning is an observed process exit. A retained process after a termination timeout is not an observed exit; `quarantined` must not be used to publish `crashed`, and a `crashed` observation must not be used to prove quarantine. U6 remains OPEN.

### 8.5 OD-1 and OD-9

OD-1's evidence-record ownership remains deferred; OD-7 does not assume an authoritative record. OD-9's approval gate is untouched. P17 remains unapproved, and this assessment does not satisfy OD-9.

### 8.6 P17 and contract rev3

P17 remains unapproved and unedited. Contract rev3 remains a frozen proposal. OD-7-A does not adopt §A3, does not amend the provider-state set, and does not authorize a public projection or diagnostic.

## 9. Remaining uncertainty

1. **Internal predicate.** The source has no quarantine marker; a timeout can leave `s.proc` installed, but it does not prove the process is alive. The exact ownership/evidence predicate remains open.
2. **Internal record and transition.** The lifecycle decision's `ProcessDisposition` names are illustrative; contract rev3's A3 is a proposal. Neither is an implementation contract.
3. **Reaping and release.** Whether a retained `Wait` owner, a second kill, or another release mechanism is correct is not established by the source or OD-7.
4. **Clearing.** Whether a later exit, successful reaping, or client replacement clears quarantine is undefined. The prior assessment's warning that a later exit must not automatically clear it remains relevant.
5. **Exposure.** Diagnostic exposure is not selected. It would require a separate approved contract; OD-7-A does not decide whether it is desirable.
6. **P17 correction.** A future P17 revision should remove `quarantined` from the required public projections or mark it non-public. No edit is made here.
7. **R11/U2/U6/U7/U3.** These remain open; OD-7-A narrows the public question but does not close them.

## 10. What this assessment does not decide

This assessment does not:

- record an OD-7 disposition, acceptance, rejection, or deferral;
- change OD-1–OD-9 or U1–U8;
- close R11, U2, U6, U7, or U3;
- adopt contract rev3 §A3;
- add, rename, or expose a public provider state;
- define the internal quarantine record, transition, retention, or reaping contract;
- edit P17, the decision register, contract rev3, or the lifecycle decision;
- authorize implementation, source changes, tests, builds, or runtime experiments.

## 11. Verification note

This is a read-only assessment. No existing artifact was modified. The decision register remains at SHA-256 `3bdc3386…`, 350 lines; OD-5 remains `47d19bf4…`; OD-6 remains `cd9d5c65…`; P17 remains `01d887c1…`; contract rev3 remains `ed0da1a7…`; the lifecycle decision remains `aa136bc9…`. No source, test, configuration, or generated file was modified; no tests, builds, benchmarks, or runtime experiments were run.

**OD-7 remains UNRESOLVED.**
