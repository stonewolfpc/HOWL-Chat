# P17 — Evidence Semantics Proposal, Revision 2 (Review Draft)

**Date:** 2026-11-09

**Status:** CORRECTED PROPOSAL AWAITING INDEPENDENT REVIEW. **Unapproved.** OD-4's orderly-shutdown outcome is accepted policy. P4-A and its non-application of retained S5 as a public-projection bar are proposed, not independently accepted policy. CP-1 remains OPEN. This is a proposed P17 revision, not a contract edit, not contract rev4, and not an implementation authorization. The original P17 proposal remains frozen and unchanged.

**Original proposal:** `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-2026-11-09.md`
**Original SHA-256:** `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf`
**Original lines:** 166

**Revision basis:** `docs/review/PLAN-howl-tts-stage3-c1-p17-revision-2026-11-09.md`, SHA-256 `56c181db3e494cd030ac6784ef7c6b803b7aac7d99a175d570d3ffd1e7a21fb7`, 374 lines; final consistency review `93a204df62470acae14393379638330f533424fcf3bb728a165b164d0d0dd873`, 125 lines. Both are plan-scoped inputs; neither validates this proposal.

**Corrected from:** 513-line candidate SHA-256 `49e9ecc2c812a361d03e739fbb571bfacc61f93fe45cffd30a8cb0bddb1055d6`; the prior correction pass closed R-1 through R-4 from the 272-line review `097fdcc20f4a6d83620c592946c36aafbb807779109f13be12331c72878283c6`.

**R-5/R-6 correction basis:** 523-line candidate SHA-256 `9cd45c468ffbaeca1512e1cfeaeba312397b7300dc7569e5b0595eae77ac3e36`; independent review `REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-2026-11-09.md`, 343 lines, SHA-256 `78fa5b7ce45825c016e9afef7298212a63d588fd069f361488ee5c901177734f`. This pass addresses R-5 and R-6 without reopening R-1 through R-4.

**N-1–N-4/N-7 edit basis:** 532-line candidate SHA-256 `8c3b7df451aee1e19c6fa4b0902b7a75f7517034b63de8adc9acd35b819e551d`; review `REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-532-2026-11-09.md`, 250 lines, SHA-256 `1910534f4b174845a41c899a670958458eb2362f9bf9cd2980e28212236d9aef`. This pass makes editorial clarifications only; it preserves R-5/R-6 and retains N-5/N-6 as OPEN scope.

**B-1 correction basis:** 536-line candidate SHA-256 `d3ca43a4feefa0c9f1108f3499b4b7796dfb0908c0acd1a0bfa6d5570dad6b40`; review `REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-536-2026-11-09.md`, 201 lines, SHA-256 `fd9e8ee98e550ea2419b9bcb0079116d2509e1783c61b5a25eed726c3a32fc78`. This pass adds only the C-B stop/OPEN condition for an uncovered separate prior public provider-state value.

**Earlier candidates:** 504-line `e1a5aaa2c68c4f30aee34440479a3d6ecd44d9220ca5cf9d187697f0734836c1`; 496-line `d783d9d8e2c5b19daadf95f9b5fce5d071434c4d1ef8c2036d6c7069361c25c3`.

**Review identity:** This corrected revision is the review candidate, not a reviewed artifact. Any qualifying review must name this exact file, revision, line count, and SHA-256; the plan-level consistency assessment cannot substitute.

**Source:** `internal/audio/tts/helper/**` at HEAD `9c7d2a1`; read-only.

**Scope:** lifecycle-affecting failure and termination observations only. Request validation, output validation, and discovery are inventoried for classification but are not lifecycle transitions.

**C1 REMAINS UNAUTHORIZED.**

No code, test, build, benchmark, runtime experiment, contract, or register file was changed. The source evidence retained below is not re-derived by this revision; E7–E9 were added from the independent source review identified in the revision plan.

## 1. Revision control and provenance

The original proposal is preserved, not overwritten. This revision changes the evidence model and projection rules only where the accepted OD-1 through OD-9 decisions require it.

| Item | Treatment |
|---|---|
| Original P17 | Frozen, unchanged; remains the author proposal of record. |
| This revision | A separate review draft; it does not approve itself or supersede the original automatically. |
| Source evidence | Original citations retained; E1–E9 retained as unresolved findings. |
| Accepted decisions | OD-1 through OD-9 applied as governing constraints. |
| U1–U8 | Remain OPEN; no revision disposition is asserted. |
| Contract rev3 | Frozen and unamended. |
| Decision register | Not edited; the OD-9 historical wording issue remains separate. |
| Technical mechanisms | Marked **PROPOSED** or **OPEN**, never as accepted policy. |

**Reading convention:** **ACCEPTED** means the owner accepted the corresponding OD decision. **RETAINED** means the item is carried forward from the original P17 or evidence inventory for provenance; it does not by itself approve or adopt the retained rule. **PROPOSED** means this revision offers a technical direction for review, not an accepted decision. **OPEN** means no determination is made. A proposal in this document is not approval of P17.

## 2. Constraints

1. Evidence is projection-supporting as well as diagnostic. It may support public provider-state projection only under an explicit, justified publication rule.
2. The supervisor is the designated authority for process-lifecycle evidence and provider-state publication. Client/session observations must be represented and retained.
3. Every relevant record binds to the process generation. Normalize relationship records by **session + supervisor generation + PID metadata**, with generation as the process-lifetime key and PID only corroborating metadata.
4. Event order and observation/processing order are distinct. Delayed delivery cannot establish causal ordering.
5. A first terminal decision is immutable once established at one linearization point. Later evidence remains retained and must not silently disappear.
6. Missing evidence is not automatically proof that an event did not occur.
7. An action is not evidence, and an exit fact is not an action.
8. Where the source cannot establish a classification, the case remains unresolved rather than being assigned a default public state.
9. No public provider state is added, removed, or renamed. The thirteen-value set remains frozen.
10. Internal evidence does not create public provider states by implication.

## 3. Model and authority

### 3.1 Record model — PROPOSED shape

```text
Observation {
  id,
  generation,
  session,
  seq_observation,
  seq_event?,          // only when the event-occurrence/acceptance point is known
  site,
  code,
  class,
  severity,
  detail,
  publication_rule?
}

Action {
  id,
  generation,
  session,
  seq_observation,
  site,
  intent,
  target,
  requested,
  accepted,
  attempted,
  completed,
  outcome,
  grace
}

ExitFact {
  id,
  generation,
  session,
  seq_observation,
  observed_exit,
  status?,
  signal?,
  wait_error?,
  observation_context
}

EvidenceLog = append-only [Observation | Action | ExitFact]
```

The exact schema, field names, allocator, retention periods, synchronization mechanism, and storage representation remain **OPEN**. The shape above is a proposal, not an accepted data contract.

`Failure` today carries `Code`, `Message`, `Retryable`, `Field`, `Detail` (`errors.go:15–21`) and no sequence, generation, class, or severity. This revision proposes adding the needed evidence fields at the observing site, not inside `Failure`; the exact placement remains **OPEN**.

### 3.2 Authority — ACCEPTED under OD-1

- One accountable authority owns evidence integrity, conflict resolution, and provider-state publication.
- The supervisor is the designated authority for process-lifecycle evidence and provider-state publication.
- Client/session observations are represented and retained; supervisor authority does not authorize silently discarding them.
- Only sufficiently justified evidence may support a public provider-state projection.
- A diagnostic label, by itself, selects no provider state. A projection requires an explicit publication rule.

The record schema and publication implementation remain **OPEN**.

### 3.3 Ordering — ACCEPTED boundary, mechanism OPEN

A sequence number allocated at observation or processing time establishes **observation order only**. It does not establish the order in which the underlying events occurred.

Event order requires evidence tied to the event occurrence or acceptance point. Wall-clock timestamps alone are not presumed to establish event order. The exact allocator and concurrency policy remain **OPEN** (C-10).

### 3.4 Terminal publication — ACCEPTED under OD-3

1. Each process instance may acquire one immutable first terminal decision.
2. That decision is established at one well-defined linearization point.
3. Later observations cannot overwrite or replace the established terminal decision.
4. Later observations remain retainable evidence associated with the same process instance.
5. Late evidence must not silently disappear merely because terminal publication occurred.
6. Evidence or operations associated with an old instance must never modify a replacement instance.

The exact synchronization mechanism, the complete terminal set, the late-evidence buffering/reconciliation design, and whether later evidence may affect diagnostics or non-terminal projections remain **OPEN**. Terminal immutability does not decide those questions.

### 3.5 Public projection eligibility — ACCEPTED boundary

A public projection is permitted only when all of the following hold:

- the evidence is sufficiently justified;
- the record is bound to the relevant generation;
- the publication rule is explicit;
- the target state is in the frozen thirteen-value provider-state set; and
- no terminal latch or generation boundary forbids the write.

Internal classifications such as `outcome_unknown`, `terminated_by_action`, and `quarantined` are not public provider states.

## 4. Evidence classes

**Status:** RETAINED source classification from original P17 §3. Inclusion in this revision is not owner acceptance, and “Fatal capable” is not newly accepted policy.

| Class | Definition | Fatal capable | Status |
|---|---|---|---|
| handshake | pre-hello negotiation or hello validation | yes | RETAINED |
| protocol | frame, envelope, or response rule violation after transport establishment | yes | RETAINED |
| transport | write failure, stream closure, receive failure | yes | RETAINED |
| exit | result of `Wait` | fact, not severity | RETAINED |
| timeout | request, shutdown, or reaping deadline expiry | action outcome | RETAINED |
| request | validation, busy, cancellation | no | RETAINED |
| output | artifact validation after synthesis | no, but may trigger cleanup | RETAINED |
| discovery | helper path or presence | no | RETAINED |

“Fatal capable” means the evidence may require lifecycle action. It does **not** establish a process exit, a public state, or causation. This is retained source classification, not a new policy determination.

## 5. Severity and publication rules

- **Fatal**: the observation invalidates the session or requires termination; no further request can succeed on that instance. Fatal does not mean “crashed.”
- **Diagnostic**: recorded for diagnosis and retained. It may support a projection only under an explicit publication rule; the label alone selects nothing.
- **Exit fact**: records what happened, never a severity.
- **Action outcome**: records what was requested and whether it completed.

| # | Rule | Status |
|---|---|---|
| S1 | Severity is bound at the site; `timeout` is diagnostic at `client.go:439` and an action outcome at `supervisor.go:285`. | RETAINED |
| S2 | One code may appear at more than one severity. | RETAINED |
| S3 | Fatal is not synonymous with crash: handshake fatal ≠ crash. | RETAINED |
| S4 | Diagnostic and projection roles are separate but not mutually exclusive. Evidence may be retained diagnostically and may support a provider-state projection only under an explicit, justified publication rule; a diagnostic label alone does not select a state. | ACCEPTED (OD-1) |
| S5 | an action never overrides a pre-existing fatal observation | RETAINED |
| S5a | An action never overrides a pre-existing terminal decision. Action causation requires OD-8’s A–E standard; action-before-exit alone is insufficient. | ACCEPTED (OD-3, OD-8) |
| S5-scope | The scope of retained S5 beyond the proposition that an action never overrides a pre-existing fatal observation is unresolved; no broader reading is accepted, inferred, or applied to OD-4 or public projection. | OPEN (tracking only; not a gate) |
| S6 | Provider state is not derived from disposition or from an unobserved outcome. In particular, no `crashed` projection follows from a fatal protocol/transport observation alone. | RETAINED |
| S7 | A confirmed process exit without established orderly shutdown projects `unavailable` unless the proposed P4-A preserves an independently justified prior **public** failure value. If an established immutable terminal exists, Case A controls and this no-orderly `unavailable` fallback does not apply; the terminal evidence remains immutable, while Case A's active-terminal and CP-1 guards govern the public projection. If such a prior public value exists but its post-exit compatibility or applicability is unresolved, that branch is **OPEN** rather than `unavailable`. An established prior terminal is not replaced by later exit evidence, but that immutability does not by itself preserve an active or `crashed` terminal value; a justified internal classification does not by itself establish a prior public state. P4-A and P4 Cases A–D govern the prior-failure/publication boundary. A prior public `crashed` value that would otherwise be preserved is governed by CP-1 OPEN. | ACCEPTED (OD-3, OD-4) boundary; P4-A PROPOSED / OPEN; CP-1 OPEN |
| S8 | An independently established orderly shutdown projects `absent` only when the proposed shared P4-A eligibility predicate is satisfied. If `absent` is excluded, preserve an established terminal only through Case A's CP-1 and active-terminal guards, or preserve an applicable prior public failure where the proposed P4-A permits; otherwise the public outcome remains **OPEN**. No active operational state survives a confirmed exit. **ACCEPTED (OD-4):** a bare fatal observation is not itself a failure classification, a public failure, a terminal, or a bar to `absent`; it does not block OD-4's orderly-shutdown projection. **PROPOSED, not independently accepted:** retained S5 is not applied here as a public-projection restriction; its broader scope is tracked under S5-scope. A genuinely unresolved failure classification (establishment or publication), publication value, or prior-public-failure applicability remains OPEN under P4-A. A prior public `crashed` value that would otherwise be preserved is governed by CP-1 OPEN. | ACCEPTED (OD-3, OD-4): orderly-shutdown and terminal boundaries; P4-A PROPOSED / OPEN; S5 non-application PROPOSED; S5-scope OPEN; CP-1 OPEN |
| S9 | `crashed` is the OD-5 internal classification for an observed process exit with no established orderly shutdown; it is not an `ExitFact` field, does not independently select a public state, and does not replace a prior terminal. A live process must not be classified `crashed`. | ACCEPTED (OD-5) |
| S10 | `outcome_unknown` is internal only; it is not a public provider state. | ACCEPTED (OD-4) |
| S11 | `terminated_by_action` is internal only and has no public projection; deliberate-action causation requires OD-8’s A–E evidence. When recorded, retain the C-probable/C-established qualifier alongside it; do not reduce it to a bare boolean or bare label. | ACCEPTED (OD-4, OD-8) |
| S12 | `quarantined` is internal-only supervisor ownership/reaping uncertainty for the exact generation; it is not a provider state and has no public projection. | ACCEPTED (OD-7) |
| S13 | OD-4 controls the public post-exit projection. OD-5’s `crashed` classification does not create a second public route; no rule in this revision selects public `crashed`. Selection is distinct from prior-value preservation under CP-1. | ACCEPTED (OD-4, OD-5) |
| S14 | Whether a previously published public `crashed` value may be preserved under OD-4 or must remain historical evidence under OD-5 is not settled by this revision. It is tracked as **CP-1 OPEN** and is not decided by implication. | OPEN (governance interaction) |

## 6. Observation inventory

The original inventory and source citations are retained. The final column applies the accepted publication boundaries; it does not decide U1–U8.

| ID | Site | Code | Class | Severity | Current handling | Revision treatment | Status |
|---|---|---|---|---|---|---|---|
| O1 | `discovery.go:39,42,61` | helper_missing | discovery | diagnostic | `setUnavailable` (`supervisor.go:132`) | Diagnostic evidence retained; the diagnostic label alone cannot select `unavailable`. Projection rule remains **OPEN**. | RETAINED / OPEN |
| O2 | `supervisor.go:147` | helper_start_failed | discovery | diagnostic | `setUnavailable`, no process | Same boundary as O1: no projection from the diagnostic label alone; rule remains **OPEN**. | RETAINED / OPEN |
| O3 | `supervisor.go:185–190` | helper_crashed | exit | fact | provider crashed, process exited | OD-5 classifies this internally as `crashed` only if the exit was observed and orderly shutdown was not established; `crashed` is not an `ExitFact` field. Route in order: (A) established terminal → route through Case A, including CP-1 and active-terminal guards; (B) pre-existing fatal observation without terminal → P4 Case B; (C) non-fatal prior failure without fatal observation or terminal → P4 Case C; (D) no prior terminal/fatal/failure → P4-D. O3 is an observed exit with no established orderly shutdown, so `absent` is excluded. If an established immutable terminal exists, Case A controls and this no-orderly fallback does not apply; apply Case A's active-terminal and CP-1 guards. Otherwise, preserve an established, applicable prior public failure; if one exists but compatibility/applicability is unresolved, the branch is **OPEN**; otherwise `unavailable` applies. An internal-only classification does not establish a prior public state. Do not publish a new `crashed` value; a prior public `crashed` value that would otherwise be preserved is governed by CP-1 OPEN. U4 remains **OPEN**. | RETAINED / ACCEPTED (OD-3, OD-4, OD-5); U4 OPEN; P4-A PROPOSED / OPEN; CP-1 OPEN |
| O4 | `supervisor.go:178–180` | none | exit | fact | absent, provider absent | Apply the proposed shared P4-A predicate after confirmed exit. `absent` only when orderly shutdown is established, no immutable terminal exists, no prior failure classification remains unresolved as to establishment, no established prior failure classification has an unresolved publication status or public value, and no applicable prior public failure exists; active operational states are not applicable, and a prior public `crashed` value that would otherwise be preserved is governed by CP-1 OPEN. A bare fatal observation is not automatically a failure classification, public failure, or terminal and does not block `absent`; retained S5 is not applied as a public-projection restriction (proposed, not independently accepted; S5-scope OPEN). Unresolved classification establishment, publication, or public value, or prior-public-failure applicability, is an OPEN branch. If `absent` is excluded, route an established terminal through Case A's CP-1 and active-terminal guards, or preserve a compatible prior public failure where the proposed P4-A permits; otherwise the branch is **OPEN**. No active operational state survives exit, and `degraded` is not presumed compatible. | RETAINED / ACCEPTED (OD-3, OD-4) boundaries; P4-A PROPOSED / OPEN; S5 non-application PROPOSED; S5-scope OPEN; CP-1 OPEN |
| O5 | `supervisor.go:179–181` | none, err retained in `waitErr` | exit | fact | provider absent despite nonzero exit | Route in order: (A) established terminal → route through Case A, including CP-1 and active-terminal guards; (B) fatal observation without terminal → P4 Case B; (C) non-fatal prior failure without fatal observation or terminal → P4 Case C; (D) no prior terminal/fatal/failure → P4-D. A nonzero exit does not establish orderly shutdown, so `absent` is excluded. If an established immutable terminal exists, Case A controls and this no-orderly fallback does not apply; apply Case A's active-terminal and CP-1 guards. Otherwise, preserve an established, applicable prior public failure; if one exists but compatibility/applicability is unresolved, the branch is **OPEN**; otherwise `unavailable` applies. An internal-only classification does not establish prior publication. A prior public `crashed` value that would otherwise be preserved is governed by CP-1 OPEN. `crashed` remains an internal classification, not an exit fact. U5 remains **OPEN**. | RETAINED / OPEN (U5); P4-A PROPOSED / OPEN; CP-1 OPEN |
| O6 | `client.go:299–303` | internal_error → helper_crashed | transport | fatal | failAll; no kill for this code | Fatal transport evidence is retained; no exit is established. Preserve the last justified current state; do not publish `absent` or `crashed`. If no current state is justified, no projection is selected and the case remains **OPEN** (C-01). | RETAINED / ACCEPTED (OD-4); C-01 OPEN |
| O7 | `client.go:304–306` | invalid_frame, frame_too_large, protocol_violation | protocol/transport | fatal | failAll + Kill | No `crashed` from the fatal observation; termination does not establish exit or causation. Preserve the last justified current state; do not publish `absent`. If no current state is justified, no projection is selected and the case remains **OPEN** (C-01). | RETAINED / ACCEPTED (OD-4); C-01 OPEN |
| O8 | `client.go:309–313` | decode errors | protocol | fatal | failAll + Kill | Same boundary as O7: no exit established; preserve the last justified current state; do not infer `absent` or `crashed`; no-state case **OPEN** (C-01). | RETAINED / ACCEPTED (OD-4); C-01 OPEN |
| O9 | `client.go:323` | protocol_violation | protocol | fatal | protocolViolation + Kill | Same boundary as O7: no exit established; preserve the last justified current state; no `absent` or `crashed`; no-state case **OPEN** (C-01). U4 remains **OPEN**. | RETAINED / ACCEPTED (OD-4); U4 OPEN; C-01 OPEN |
| O10 | `client.go:331` | protocol_violation | protocol | fatal | protocolViolation + Kill | Same boundary as O7: no exit established; preserve the last justified current state; no `absent` or `crashed`; no-state case **OPEN** (C-01). U4 remains **OPEN**. | RETAINED / ACCEPTED (OD-4); U4 OPEN; C-01 OPEN |
| O11 | `client.go:345` | protocol_violation | protocol | fatal | protocolViolation + Kill | Same boundary as O7: no exit established; preserve the last justified current state; no `absent` or `crashed`; no-state case **OPEN** (C-01). U4 remains **OPEN**. | RETAINED / ACCEPTED (OD-4); U4 OPEN; C-01 OPEN |
| O12 | `client.go:421` | helper_crashed | transport | fatal | returned only; no termination | No exit established; preserve the last justified current state; no `absent` or `crashed`. If no current state is justified, no projection is selected and the case remains **OPEN** (C-01). | RETAINED / ACCEPTED (OD-4); C-01 OPEN |
| O13 | `client.go:439` | timeout | timeout | diagnostic | request fails | No lifecycle projection from a request timeout alone. | RETAINED |
| O14 | `client.go:441` | cancelled | request | diagnostic | request fails | No lifecycle projection from cancellation alone. | RETAINED |
| O15 | `client.go:248–253` | incompatible, minor_unsupported, others | handshake | fatal | Terminate then publish incompatible or unavailable | Apply P2’s exact categories; no generic default. Path-by-path classification remains **OPEN** (U8). | RETAINED / OPEN (U8) |
| O16 | `protocol.go:502–545` | protocol_violation | handshake | fatal | falls to `unavailable` at `client.go:253` | Apply P2 exactly: missing mandatory capability → `incompatible`; duplicate identities, invalid semantic descriptors, and malformed helper version/limits mismatch → `unavailable`; misplaced descriptors, unsupported extension versions, unknown optional/provider feature requiredness, and unlisted classifications remain unclassified. Do not generalize all malformed or ambiguous cases to `unavailable`. U8 remains **OPEN**. | RETAINED / OPEN (U8) |
| O17 | `client.go:578–594` | wire/protocol failures | load/transport | diagnostic | provider unavailable | Load-path wire/protocol failure; this is not a hello/handshake compatibility event. Do not apply OD-6 or P2 here. No exit is established; preserve the last justified current state; no `absent` or `crashed`. If no current state is justified, no projection is selected and the case remains **OPEN** (C-01); a separate load-path rule remains **OPEN**. | RETAINED / ACCEPTED (OD-4); load-path rule OPEN; C-01 OPEN |
| O18 | `client.go:756` | helper_crashed | transport | fatal | provider crashed while process may be alive | A live process must not be `crashed`; retain the evidence. No exit is established; preserve the last justified current state; no `absent`. If no current state is justified, no projection is selected and the case remains **OPEN** (U6/C-01). | RETAINED / ACCEPTED (OD-4); U6 OPEN; C-01 OPEN |
| O19 | `client.go:771,783,793,796` | timeout, cancelled, degraded | request | diagnostic | provider ready or degraded | The diagnostic label alone selects no state; an explicit publication rule is required. Current handling is not accepted as a rule. | RETAINED / OPEN |
| O20 | `supervisor.go:283–286` | timeout | timeout | action outcome | returned; process state unchanged | Action outcome only; no terminal projection. | RETAINED |

## 7. Termination actions

The original action inventory is retained. Action identity and outcome must be recorded; an action is still not proof of causation.

| ID | Site | Intent (proposed) | Grace | Intent recorded today | Revision boundary | Status |
|---|---|---|---|---|---|---|
| T1 | `client.go:248` | handshakeFailure | ShutdownTimeout | no | Retain action identity and outcome; no public state by itself. | RETAINED / PROPOSED intent |
| T2 | `client.go:305` | protocol | default | no | `Kill` aliasing obscures intent (E4); retain distinct action evidence. | RETAINED / PROPOSED intent |
| T3 | `client.go:312` | protocol | default | no | Same as T2. | RETAINED / PROPOSED intent |
| T4 | `client.go:366` | protocol | default | no | Same as T2. | RETAINED / PROPOSED intent |
| T5 | `client.go:654` | shutdown | ShutdownTimeout | no | Retain completion/outcome; orderly shutdown must be established, not assumed. | RETAINED / PROPOSED intent |
| T6 | `client.go:661` | shutdown | ShutdownTimeout | no | Same as T5. | RETAINED / PROPOSED intent |
| T7 | `client.go:664` | shutdown | ShutdownTimeout | no | Same as T5. | RETAINED / PROPOSED intent |
| T8 | `client.go:668` | shutdown | ShutdownTimeout | no | Same as T5. | RETAINED / PROPOSED intent |
| T9 | `client.go:861` | cleanup | ShutdownTimeout | no | Cleanup errors must be retained (E8); no success assumed. | RETAINED / PROPOSED intent |
| T10 | `supervisor.go:275,278` | shutdown-grace-expiry | remaining grace | no | Timeout is an action outcome, not proof of exit, quarantine, or successful reaping. | RETAINED / OPEN |

`Kill` is an alias for `Terminate` with the default grace (`supervisor.go:291–293`), so T2–T4 are indistinguishable from T10 in evidence. Deliberate-action causation requires OD-8’s A–E standard; U5 remains **OPEN**. No termination action has a public provider state of its own.

## 8. Publication rules

These rules replace the original precedence table. They define eligibility for a projection; they do not decide the open technical mechanisms.

### P1 — Authority and evidence — **ACCEPTED (OD-1)**

- The supervisor is the designated publication authority for process-lifecycle evidence.
- Client/session observations are retained and cannot be silently discarded.
- A projection must cite sufficient, generation-bound evidence and an explicit rule.
- Internal evidence does not create a new public state.

### P2 — Handshake compatibility — **ACCEPTED (OD-6); path classification OPEN (U8)**

1. Missing mandatory service capability → `incompatible`.
2. Unknown required core capability → `incompatible`.
3. Known required capability at an unsupported required version → `incompatible`.
4. Duplicate identities → `unavailable`.
5. Invalid semantic descriptors → `unavailable`.
6. Malformed helper version or limits mismatch → `unavailable`.
7. Misplaced descriptors, unsupported extension versions, unknown optional/provider feature requiredness, and any unlisted classification remain **unclassified**.

The exact path-by-path classification of the nine `validateHello` exits remains **OPEN** (U8). A failure-code lookup does not replace the observed condition and category.

### P3 — Fatal protocol or transport observation — **ACCEPTED (OD-4, OD-5); no-justified-state case OPEN**

A fatal protocol or transport observation is retained as internal evidence. It does **not** establish a process exit, a terminal decision, or a public state.

- Do not publish `crashed` from the observation alone.
- Do not infer or publish `absent`.
- No observed exit is not evidence establishing process absence.
- If no exit is established, preserve the last justified current state. A state that was merely displayed, stale, or not independently justified is not made current by this rule.
- If no current provider state has been justified, select no projection and keep that precise case **OPEN** (C-01).
- If an exit is later observed, apply P4 and the terminal latch.

### P4 — Confirmed exit, prior evidence, and routing — **ACCEPTED (OD-3, OD-4, OD-5) boundaries; P4-A PROPOSED / OPEN**

A confirmed exit is a fact distinct from orderly shutdown, causation, internal classification, and publication. `crashed` is an internal classification derived from an observed exit when orderly shutdown is not established; it is not an `ExitFact` field and not a public state by itself.

**Predicates.** Keep these distinct; do not infer one from another:

1. An **internal fatal observation** exists in the evidence record.
2. A **failure classification** is independently justified by that observation or other evidence.
3. A **previously published public provider-state value** is established by evidence of the publication.
4. An **established immutable first terminal decision** exists under OD-3.

Observed process exit and independently established orderly shutdown are separate lifecycle facts, not one of those four predicates. An internal fatal observation does not establish a failure classification, a public state, or a terminal. A justified failure classification does not establish that it was published or became the immutable first terminal. A previously published public state is not interchangeable with an internal-only classification. Do not claim that a public state can be preserved or overwritten unless the relevant prior publication and evidence are established.

**Failure-classification threshold — OPEN.** Observing a fatal event, establishing an internal failure classification from it, publishing a public failure value, and establishing a terminal are separate determinations. The available accepted decisions do not adopt one universal evidence threshold for converting a fatal observation into a failure classification. Where the event's accepted semantics and available evidence already establish the classification, this proposal does not require additional corroboration; where they do not, the threshold remains **OPEN** and the classification must not be assumed. This is an evidence-specific determination, not a license to dismiss the observation.

**Classification establishment — three-way determination.** For every prior failure classification advanced by the record, the establishment determination is exactly one of: **established** (the threshold rule is settled and the available evidence satisfies it); **rejected** (the threshold rule is settled and the evidence fails it, or the accepted semantics exclude the classification); or **unresolved** (the threshold rule, its application to the evidence, or the classification determination itself is unresolved). A settled threshold rule alone is not establishment. **Unresolved establishment blocks `absent`.**

**P4-A — Shared `absent` eligibility predicate — PROPOSED, not accepted policy.** After a confirmed process exit, `absent` is eligible only when **all** of the following gates hold:

1. orderly shutdown is independently established;
2. no established immutable first terminal decision exists for the same process generation;
3. no prior failure classification remains **unresolved as to establishment**. A settled threshold rule alone does not satisfy this gate; the evidence application must establish or reject the classification;
4. no **established** prior failure classification has an unresolved publication status or unresolved public value. Publication status is exactly `published` (with the public value identified), `excluded-internal-only` (positively established under the applicable publication rule or other evidence), or `unresolved`. Mere absence of publication evidence is `unresolved`; an internal-only classification does not establish publication exclusion by itself;
5. no independently justified prior **public** failure value is established as applicable after exit or remains unresolved on post-exit compatibility or applicability. An active operational state, including active `degraded`, is not applicable and cannot be preserved; a separately asserted post-exit interpretation of that state remains subject to C-C and, while unresolved, prevents `absent`. A prior public `crashed` value that would otherwise be compatible and applicable is governed by **CP-1 OPEN** and holds the branch **OPEN**; a `crashed` value established as incompatible is historical evidence and does not fail this gate.

**Public-effect dimensions.** For a prior failure classification, resolve (a) establishment: established, rejected, or unresolved; (b) if established, publication status: published with the public value identified, excluded-internal-only, or unresolved; (c) post-exit compatibility; and (d) post-exit applicability. Gate 3 covers (a); gate 4 covers (b); gate 5 covers (c) and (d). A threshold rule and its evidence application are components of (a); a settled rule alone is not establishment. Publication exclusion is a positive determination, not the absence of publication evidence. Proof that a value was published establishes publication status; it does not by itself establish that the value is independently justified as current projection evidence. These dimensions are distinct; the gates are not interchangeable. Compatibility does not establish applicability; a compatible but non-applicable value does not fail gate 5.

A bare fatal observation is not a failure classification, a public failure, a terminal, or prior failure evidence, and it does not block `absent`. Gate 3 fails when the record has advanced a prior failure classification whose establishment remains unresolved. Gate 4 fails when an established prior failure classification has an unresolved publication status or public value, including where publication evidence is merely missing. A bare observation alone does not advance a classification, fails neither gate, and does not invoke retained S5; this non-application is part of the proposed P4-A, and `S5-scope` remains **OPEN**. The same establishment/threshold/publication/public-value resolution standard applies to fatal and non-fatal classifications; the difference is routing, not the standard. With orderly shutdown and no applicable or unresolved prior public failure, OD-4's `absent` row applies. When gates 3 and 4 are resolved, an applicable failure classification is routed through gate 5 and the applicable Case B or C branch. An established internal-only classification is not a public value only when its publication status is positively resolved as `excluded-internal-only`; the classification itself does not establish that exclusion. A merely asserted classification, an unresolved establishment, or an unresolved publication status cannot satisfy the gates by silence. A public failure value must have evidence of publication, and its post-exit compatibility and applicability are separate questions. A value that is established as incompatible does not defeat gate 5; one that is compatible and applicable, or unresolved on either question, does. `degraded` is not presumed post-exit-compatible merely because it is failure-related. No active operational state is eligible for preservation after confirmed exit.

**Status of P4-A.** P4-A is a proposed routing restatement, not an accepted owner decision. Gates 1 and 2 restate OD-3/OD-4 boundaries; gate 5 restates OD-4's prior-failure and public-publication boundary; gate 3 requires a resolved classification establishment; gate 4 requires resolved publication status and public value for an established classification. Retained S5 is not a gate, but that non-application is part of P4-A's proposal, not independently accepted policy. S5-scope remains OPEN.

**Projection from P4-A.**

**Terminal precedence:** Whenever Gate 2 fails because an established immutable terminal exists, Case A controls every projection branch below; no `absent` or `unavailable` fallback overrides the active-terminal or CP-1 guards.

- **All five gates hold:** project `absent`.
- **Gate 2 fails:** Case A controls; apply the Case A terminal guards for CP-1 and active operational states before any public projection.
- **Gate 3 fails, with no applicable gate-5 failure:** `absent` is not eligible. With orderly shutdown established, the outcome is **OPEN** because the prior failure classification's establishment remains unresolved. Without orderly shutdown, OD-4's `unavailable` fallback applies unless an established terminal exists, in which case Case A controls under the terminal-precedence rule above.
- **Gate 4 fails, with no applicable gate-5 failure:** `absent` is not eligible. With orderly shutdown established, the established classification's publication status or public value remains unresolved, so the outcome is **OPEN**. Without orderly shutdown, OD-4's `unavailable` fallback applies unless an established terminal exists, in which case Case A controls under the terminal-precedence rule above.
- **Gate 5 fails:** `absent` is not eligible. Preserve the prior public failure where its post-exit compatibility and applicability are established and it is not an active operational state; if either remains unresolved, the outcome is **OPEN**. If the prior public value is `crashed` and would otherwise be compatible and applicable, apply **CP-1 OPEN**; retain it as historical evidence and do not decide preservation or exclusion here. A value established as incompatible does not fail gate 5 and does not compete with the exit-state projection. Do not silently replace an applicable or unresolved public failure with `absent` or `unavailable`.
- **No orderly shutdown:** `absent` is not eligible. If an established immutable terminal exists, Case A controls and this `unavailable` fallback does not apply; apply Case A's active-terminal and CP-1 guards. Otherwise, apply OD-4's `unavailable` fallback unless an applicable prior public failure is preserved. If a prior public failure exists but its compatibility or applicability is unresolved, that branch is **OPEN**. A prior public `crashed` value that would otherwise be preserved is governed by **CP-1 OPEN**. A nonzero exit alone does not establish orderly shutdown.

### Case A — Established prior terminal decision — **OD-3/S5a terminal immutability ACCEPTED; active-state exclusion ACCEPTED (OD-4); terminal-set definition OPEN; CP-1 OPEN**

OD-3 and S5a govern the immutable terminal decision. This section does not define the terminal set; that remains **OPEN**.

- **No exit:** preserve the terminal; later evidence is retained but cannot replace it.
- **Confirmed exit — active operational terminal guard:** Before preserving a terminal value after exit, determine whether that public value is an active operational state, including `ready`, active `degraded`, or any state whose meaning asserts current operation. Such a value is not preservable. Retain the immutable decision and later evidence, do not project the active state, and keep the public projection **OPEN** because the terminal decision still exists; this proposal does not invent a replacement terminal projection. If it cannot be determined whether the value is active, it is not preserved. This guard applies through Gate 2 and B-D and does not resolve the terminal set.
- **Confirmed exit — CP-1 terminal guard:** If the terminal's public value is `crashed`, **CP-1 OPEN** applies. Retain the immutable decision and later evidence; do not assert preservation, and do not assert historical-only exclusion. The public projection remains **OPEN** pending CP-1. This guard applies through Gate 2 and B-D.
- **Confirmed exit; no terminal guard applies; orderly shutdown established:** preserve the terminal. Later orderly-shutdown evidence is retained but must not overwrite or replace the immutable decision. Do not project `absent` over it.
- **Confirmed exit; no terminal guard applies; orderly shutdown not established:** preserve the terminal. A `crashed` classification may be recorded only as internal later evidence if OD-5 holds; it does not replace the terminal or become public.
- **Public projection:** follows the established terminal and the accepted projection rules, subject to the active-terminal and CP-1 guards. This revision does not invent a new terminal projection rule.

### Case B — Pre-existing fatal observation, no terminal decision — **OD-4 outcomes ACCEPTED; P4-A PROPOSED; S5 non-application PROPOSED; S5-scope OPEN; unresolved classifications OPEN**

S5 is **RETAINED** verbatim, not accepted policy. OD-3/S5a do not apply because no immutable terminal exists. S5's retained wording concerns an action overriding a fatal observation; it does not establish that every fatal observation prohibits `absent`. The non-application of S5 as a public-projection bar is **PROPOSED**, not independently accepted policy. Its broader scope is tracked as **S5-scope OPEN**. No broader reading is applied as a gate here.

A bare fatal observation is not a failure classification, a public failure, a terminal, or prior failure evidence. With orderly shutdown and no justified prior failure or prior public failure, OD-4's `absent` row applies. A failure classification advanced by the record with an unresolved threshold, establishment, publication, or public value remains **OPEN** under B-B. A merely asserted internal-only classification is not established and remains under Gate 3; do not route it as C-B.

**B-A — Bare fatal observation only.** B-A applies only when no prior failure classification has been advanced. If a classification has been advanced, B-B controls even when its establishment is unresolved or the classification is rejected.

- **No exit:** apply P3. Preserve the last justified current state; do not project `absent` or `crashed`. If no current state is justified, no projection and **OPEN** (C-01).
- **Exit observed; orderly shutdown independently established:** no prior failure classification has been advanced, and no independently established prior public failure value, terminal, or accepted public-projection restriction controls. Project `absent` under OD-4. The bare observation does not invoke retained S5; this non-application is **PROPOSED**, and `S5-scope` remains **OPEN**.
- **Exit observed; orderly shutdown not established:** `unavailable` under OD-4 unless an applicable prior public failure is preserved. If a prior public failure exists but compatibility or applicability is unresolved, that branch is **OPEN**. `crashed` may be recorded only as an internal classification under OD-5. Do not publish a new `crashed` value; a prior public `crashed` value that would otherwise be preserved is governed by CP-1 OPEN.

**B-B — Fatal observation with a prior failure classification advanced or established.** Retain the observation and classification separately. Establish the classification before routing; a settled threshold rule alone does not do so. The classification does not establish a public value or a terminal by implication.

- **No exit:** apply P3.
- **Exit observed; orderly shutdown independently established; establishment unresolved:** the outcome is **OPEN**. `absent` is not eligible. The threshold rule may be settled while its application to the evidence remains unresolved; that state is still unresolved for Gate 3.
- **Exit observed; orderly shutdown independently established; classification rejected:** the classification does not qualify as a prior failure. Retain it as historical evidence. If no independently established prior public value exists and the other gates hold, OD-4's `absent` row applies; any independently established prior public value remains governed by B-C and Gate 5.
- **Exit observed; orderly shutdown independently established; established, publication status unresolved:** the outcome is **OPEN**. Missing publication evidence is not exclusion, and `absent` is not eligible.
- **Exit observed; orderly shutdown independently established; established, publication excluded-internal-only:** the classification remains non-public. If a separate prior public value exists, B-C/Gate 5 governs it (or C-C if it is an active operational state), including CP-1 where applicable; if its compatibility or applicability is unresolved, the branch is **OPEN**; `absent` is not eligible for that value. With no separate independently established prior public value and the other gates satisfied, OD-4's `absent` row applies.
- **Exit observed; orderly shutdown independently established; established, publication status published:** apply B-C, including its CP-1 and active-state handling.
- **Exit observed; orderly shutdown not established:** `unavailable` under OD-4 unless an applicable prior public failure is preserved. If one exists but compatibility or applicability is unresolved, that branch is **OPEN**. An internal-only classification does not establish public preservation. Do not publish a new `crashed` value; a prior public `crashed` value that would otherwise be preserved is governed by CP-1 OPEN.

**B-C — Fatal observation with a previously published public failure value.** Preserve the value under OD-4 only where it remains independently justified, post-exit-compatible, and applicable. Compatibility does not establish applicability; a compatible but non-applicable value is not preserved and follows ordinary OD-4 exit handling. An active operational state, including active `degraded`, is not preservable under this rule; C-C controls. If the value is `crashed` and would otherwise be compatible and applicable, do not decide preservation here: **CP-1 is OPEN**; retain the value as historical evidence and do not treat it as preserved or excluded. If compatibility or applicability is unresolved, the branch is **OPEN**; do not substitute `absent` or `unavailable`. If the value is established as incompatible, retain it only as historical evidence and do not preserve it as the current public state; it is not an active operational state. After an incompatible value is reduced to history, with orderly shutdown and no other controlling state, OD-4's `absent` row applies. Without orderly shutdown, `unavailable` applies. Retained S5 is not applied as a public-projection bar (**PROPOSED**; `S5-scope` **OPEN**).

**B-D — Fatal observation with an established immutable terminal decision.** Route through Case A, including its active-terminal and CP-1 guards, and apply OD-3/S5a. Later evidence cannot overwrite the terminal; the guards do not decide CP-1 or define the terminal set.

### Case C — Non-fatal prior failure, no fatal observation and no terminal — **P4-A PROPOSED; Case C subcategories; unresolved interactions OPEN**

OD-3 and S5a do not apply because no terminal exists; S5 does not apply because there is no pre-existing fatal observation, and `S5-scope` remains **OPEN**. OD-4 and P4-A remain the public projection policy. The three subcategories below are exhaustive for this case. A prior classification that remains unresolved as to establishment fails Gate 3 and remains **OPEN**; absent an independently established prior public value, it must not be routed through C-A, C-B, or C-C to `absent`.

**C-A — Previously published failure state that is not an active operational state.** If the prior public failure value is independently justified and its post-exit compatibility and applicability are established, preserve it under OD-4. If the value is `crashed` and would otherwise be compatible and applicable, apply **CP-1 OPEN**; retain it as historical evidence and do not treat it as preserved or excluded. `absent` is excluded by P4-A gate 5. If compatibility or applicability remains unresolved, `absent` is excluded and the outcome is **OPEN**. Compatibility alone does not establish applicability; a compatible but non-applicable value is not preserved and does not fail gate 5. An active `degraded` state is not C-A; it belongs in C-C even when it is failure-related.

If the prior public failure value is independently established as **not post-exit-compatible**, retain it as historical evidence, do not preserve it as the current public provider state, and do not reclassify it as an active operational state. Route the ordinary confirmed-exit question under P4-D: with orderly shutdown independently established and no other controlling prior state, project `absent`; without orderly shutdown, project `unavailable`; if applicability or a competing accepted constraint remains unresolved, retain **OPEN** with that exact reason. Do not erase the historical existence of the incompatible failure value.

**C-B — Established internal-only prior failure classification.** Retain the evidence and classification. This subcategory applies only when the classification is established and publication status is positively resolved as `excluded-internal-only`; a merely asserted internal-only classification, an unresolved establishment, or missing publication evidence is not C-B and remains **OPEN** under Gate 3 or Gate 4 as applicable. Do not infer publication exclusion from the internal-only label or from silence. Do not infer that it was previously published and do not project it as a public provider state. The same establishment/threshold/publication/public-value standard as B-B applies.

Before applying the internal-only branch, distinguish this classification from any separate prior public provider-state value. An internal-only classification does not establish that no separate prior public value exists. C-B applies only when no fatal observation and no terminal exist; if either exists, Case B or Case A controls instead. A separate prior public provider-state value is not governed by C-B's `absent` or `unavailable` fallback. Route a non-active prior public failure through C-A and an active operational state through C-C. If a separate prior public provider-state value is present and is neither a prior public failure value nor an active operational state, C-B does not proceed to the ordinary `absent`/`unavailable` outcome; the combination remains OPEN. Preserve a prior public failure value only where it remains independently justified, post-exit-compatible, and applicable; if compatibility or applicability is unresolved, the branch is **OPEN**; if it is `crashed` and would otherwise be compatible and applicable, **CP-1 OPEN** applies and it remains historical evidence rather than a decided preservation or exclusion. If it is established as incompatible, retain it as historical evidence and apply the ordinary confirmed-exit handling under P4-D; the internal-only classification remains retained.

Only when no separate prior public value remains applicable, unresolved, or **CP-1 OPEN**, and no separate prior public value is present that is neither a prior public failure value nor an active operational state, does C-B proceed to the ordinary `absent`/`unavailable` outcome: with orderly shutdown and the other gates satisfied, `absent` is eligible; without orderly shutdown, `unavailable` applies. If the internal-only classification's establishment, threshold, publication, or public value is unresolved, `absent` is not eligible; the orderly-shutdown outcome is **OPEN**. Without orderly shutdown, `unavailable` applies only if no separate prior public value remains applicable, unresolved, or **CP-1 OPEN**, and no separate prior public value is present that is neither a prior public failure value nor an active operational state; otherwise the separate prior public value governs under the preceding paragraph, including CP-1 or an **OPEN** branch.

**C-C — Prior active operational state, including active `degraded`.** Do not preserve an active operational state after confirmed process exit. This exclusion is applied within P4-A gate 5 as an eligibility filter; it is not a separate gate and does not by itself make gate 5 fail. It does not depend on reaching this subcase. This includes `ready`, active `degraded`, and any other state whose meaning asserts current operation. A failure-related label does not make an active `degraded` state preservable. If a separately justified post-exit interpretation of a prior `degraded` value is asserted, do not assume it; require evidence, and retain the interpretation as **OPEN** while it is unresolved. Otherwise apply P4-A: `absent` only if the remaining gates are satisfied and orderly shutdown is established; without orderly shutdown, `unavailable` under OD-4. The discarded active projection is not a competing outcome.

### C-D — No prior failure (transition to P4-D)

This is not Case C. A prior classification that has been affirmatively rejected is not a prior failure; retain it as historical evidence. If no independently established prior public value remains, route ordinary confirmed exit through P4-D. Any independently established prior public value remains governed by C-A and Gate 5, including CP-1 where applicable.

### Case D (P4-D) — Ordinary confirmed exit, no applicable prior controlling state — **ACCEPTED (OD-4)**

No applicable prior terminal, fatal observation, or failure classification remains. An independently established incompatible prior public failure has been reduced to historical evidence under C-A and does not make the exit non-ordinary.

- **Exit observed; orderly shutdown independently established:** OD-4 directly projects `absent`; the proposed P4-A predicate is satisfied.
- **Exit observed; orderly shutdown not established:** project `unavailable` under OD-4. `crashed` may be recorded internally under OD-5 but is not published.

**Public `crashed`.** OD-5’s classification does not create a second public route. This revision does not select a new public `crashed` projection for any exit condition; the state remains in the frozen thirteen-value set. Whether a previously published public `crashed` value may be preserved under OD-4 or must remain historical evidence under OD-5 is **CP-1 OPEN**; this section does not decide it by implication. This includes a terminal whose public value is `crashed`; Case A routes that value to CP-1 rather than treating terminal preservation as a decision. A future new public `crashed` route requires a separate accepted decision.

A nonzero exit alone proves neither orderly shutdown nor deliberate-action causation. It must not be treated as `absent`.

### P5 — Fatal observation combined with an exit — **ACCEPTED boundary; P4-A PROPOSED; U3/U4 OPEN; C-01 OPEN**

- A fatal observation alone does not establish an exit or a terminal decision.
- If an exit is observed, P4’s combined routing governs classification and public projection; the fatal observation is retained but does not independently select `crashed`.
- If no exit is established, P3 applies: preserve the last justified current state; do not publish `absent` or `crashed`; if no current state is justified, the case remains **OPEN** (C-01).
- When an exit is established, P4-A determines `absent` eligibility; a previously published public provider state versus an internal-only classification is a required distinction.
- U3 and U4 remain **OPEN**; the no-current-state branch remains **C-01 OPEN**.

### P6 — Causation (action or fatal observation) — **ACCEPTED (OD-8); implementation OPEN**

`terminated_by_action` remains internal and has no public projection. A causal attribution requires OD-8’s A–E standard and the associated F assessment; satisfying a subset is insufficient. The exact evidence mechanism, allocator, storage, and reconciliation remain **OPEN**.

**A — Process-instance identity.** The action, action outcome, exit fact, and supporting evidence must bind to the same authoritative process instance. Use OD-2’s generation-based identity, including uniqueness across supervisor-object replacement within a session. PID alone is insufficient. Old-instance evidence must not be attributed to a replacement.

**B — Action identity and outcome.** Identify the specific deliberate action and its target. Distinguish requested, accepted, attempted, completed, rejected, timed out, and failed or otherwise unsuccessful. A request is not proof of execution, and a successful API return is not proof of causal effect.

**C — Exit fact.** Identify the actual process exit observation and preserve available status, signal, error, observation context, and instance identity. Do not substitute a generic `waitErr`, provider state, or lifecycle label for actual exit evidence.

**D — Event ordering.** Classify event order as **A-E** (established), **A-O** (observation order only), or **A-U** (unreliable or unknowable). A-E is required for any causal attribution, including **C-probable**. A-O establishes association only and cannot support C-probable or C-established. A-U cannot support causal attribution. A sequence number is authoritative only if allocated at the relevant event occurrence or acceptance point; an observation-time sequence establishes observation order only. Wall-clock timestamps alone are not presumed to establish event order.

**E — Causal mechanism.** A causal attribution requires all of:

1. a mechanism capable of producing the observed exit;
2. evidence that the mechanism actually operated on this particular process instance;
3. consistency between that mechanism and the observed exit facts;
4. an explicit assessment of material competing explanations; and
5. an explicit statement of remaining uncertainty.

A plausible mechanism narrative is insufficient. A successful action result, a kill-consistent exit status, `stopping`, or temporal correlation is not independently sufficient. Do not use the same state observation as circular proof of action identity, event ordering, and mechanism operation.

**F — Competing explanations and defeaters.** Evaluate material competing explanations using available evidence and record:

- the alternative explanation;
- evidence supporting or contradicting it;
- whether it is excluded, materially weakened, still plausible, or unresolved; and
- what evidence is missing.

Do not manufacture a defeater merely to satisfy a checklist. If no material competing explanation is identified, document the assessment scope and why none was identified. Missing evidence for a competing explanation is not evidence in favor of the preferred explanation. A material unresolved alternative limits the strength of the conclusion.

**Causal classification.** Preserve **C-unresolved** (available evidence does not support a causal conclusion), **C-probable** (A-E established, mechanism operation evidenced, exit facts consistent, evidence favors action causation over evaluated alternatives, but material uncertainty remains), and **C-established** (A-E established, strong convergent mechanism evidence, exit facts consistent, and no material competing explanation unresolved). `terminated_by_action` may be recorded only for C-probable or C-established; C-unresolved is not `terminated_by_action`. Retain the C-probable/C-established qualifier alongside the record; do not reduce it to a bare boolean or bare label. Do not assign numerical probabilities without an independently justified statistical basis. Lifecycle closure is separate from causation and must not be used as causal evidence. U5 remains **OPEN**.

### P7 — Internal classifications — **ACCEPTED (OD-4, OD-7); representation OPEN**

- `outcome_unknown` may be recorded internally when the evidence cannot support a public classification. It is not a public provider state and must not appear in the public projection column.
- `terminated_by_action` remains internal and must retain its C-probable/C-established causal qualifier. It does not replace a public state by implication.
- `quarantined` is internal-only ownership/reaping uncertainty for the exact generation. It is not a provider state, is not a fourteenth state, and has no public projection.
- The exact internal representation, transition, retention, and release rules for all three remain **OPEN**.

### P8 — Diagnostic and diagnostic-only evidence — **ACCEPTED boundary (OD-1); publication rules OPEN**

A diagnostic observation may be retained for diagnosis and may also support a projection under an explicit rule. A diagnostic label or diagnostic-only record does not select a provider state by itself.

O1, O2, and O19 currently select states despite being classified as diagnostic. The revision does not preserve that handling as policy. Either a projection rule is stated, or the observation remains non-projecting.

### P9 — Terminal latch and late evidence — **ACCEPTED (OD-3); mechanics OPEN**

- The first terminal decision for a generation is immutable.
- Later evidence is retained but cannot overwrite or republish the terminal.
- Old-generation evidence cannot modify a replacement.
- Whether later evidence may affect diagnostics or non-terminal projections remains **OPEN**.
- The exact linearization and reconciliation mechanics remain **OPEN**.

### P10 — Quarantine — **ACCEPTED (OD-7); mechanics OPEN**

Quarantine means the supervisor has explicitly recorded an unresolved ownership/reaping condition for the exact generation. It is not merely a termination timeout, not merely a claim that a process is alive, and not a public state. A later exit does not automatically establish successful reaping. The exact evidence, transition, retention, and release rules remain **OPEN**.

## 9. Evidence-loss findings

The original findings are retained. E7–E9 are added from the independent source review and remain part of the evidence history.

| ID | Finding | Source |
|---|---|---|
| E1 | `failAll` keeps only the first read failure; later failures are discarded | `client.go:369–386` |
| E2 | `SetState` overwrites `lastFailure`, so only the last published failure survives | `supervisor.go:233–241` |
| E3 | the stopping branch discards a nonzero exit from classification | `supervisor.go:176–187` |
| E4 | termination intent is not a parameter, and `Kill` aliases `Terminate` | `supervisor.go:246,291–293` |
| E5 | `waitForReader` returns nothing on timeout | `client.go:352–367` |
| E6 | stream closure is marked crashed but does not terminate for that code | `client.go:299–308` |
| E7 | late terminal evidence can be dropped by dispatch | `client.go:339–342` |
| E8 | cancellation and cleanup paths ignore errors from `Cancel` and `Terminate` | `client.go:248,654,661,664,861`; `supervisor.go:246–288` |
| E9 | `Shutdown` calls `Terminate`, then clears lifecycle and transport state without establishing successful process reaping | `client.go:654–674`; `supervisor.go:246–288` |

These findings remain unresolved. This revision does not claim they are fixed, and their presence prevents treating precedence or publication as trustworthy without separate evidence or implementation.

## 10. Unresolved cases

No classification is asserted for U1–U8. All remain **OPEN**.

| ID | Case | Why unresolved | Status |
|---|---|---|---|
| U1 | orderly EOF during shutdown vs `helper_crashed` from the reader | the reader maps generic stream errors to crash while the observer publishes absent | OPEN |
| U2 | shutdown timeout racing a protocol violation | two fatal/action outcomes with no defined winner | OPEN |
| U3 | handshake incompatible followed by a process crash | first terminal vs sequence order | OPEN |
| U4 | exit zero after a fatal protocol violation | fatal-wins rule may overstate certainty | OPEN |
| U5 | nonzero exit after a deliberate kill | crash vs terminated_by_action is not established | OPEN |
| U6 | provider `crashed` while the process is alive (`client.go:756`) | provider state conflates condition with process outcome | OPEN |
| U7 | `Terminate` timeout after the terminal projection is published | terminal immutability is undefined | OPEN |
| U8 | missing mandatory capability yields `protocol_violation` to `unavailable`, while helper-declared incompatibility yields `incompatible` | same concept, two projections | OPEN |
| CP-1 | prior public `crashed` preservation | OD-4 preserves an applicable prior failure; OD-5 defines `crashed` as observed exit without orderly shutdown. Whether a previously published public `crashed` is preservable or historical-only is not settled, including when the value is carried by an established terminal. | OPEN |

CP-1 is a governance interaction, not a U1–U8 classification case; U1–U8 remain **OPEN**.

OD-6 supplies the classification boundary for a missing mandatory capability, but it does not establish the exact source path/code classification or close U8. The accepted decisions constrain how U1–U8 may be evaluated; they do not close them.

## 11. Reviewer challenges

1. Verify every retained line and code claim against source; do not accept a classification from this draft alone.
2. Challenge each fatal assignment, especially O6, O12, O15, O16, and O18.
3. Review U1–U8 as open; any test, build, benchmark, or runtime experiment requires separate authorization and is not performed here.
4. Check that terminal immutability does not erase late evidence, that old-generation evidence cannot modify a replacement, and that every route into Case A applies the active-terminal and CP-1 guards before public preservation, including Gate 2 and B-D.
5. Confirm that no public projection uses `outcome_unknown`, `quarantined`, or `terminated_by_action`, and that the thirteen-value provider-state set is unchanged.
6. Confirm that E1–E9 remain visible and unresolved, and that the revision does not claim implementation evidence.
7. Return finding-level dispositions; do not approve P17 on the author’s or this revision’s behalf.
8. Record the exact reviewed artifact: path, revision, line count, and SHA-256. Do not treat the plan-level consistency assessment as a review of this proposal.

## 12. Change summary from the original P17

| Change | Original issue | Revision treatment |
|---|---|---|
| C-01 | Fatal protocol/transport projected `crashed` without an observed exit. | Removed; fatal evidence retained, exit not inferred. With no exit, preserve the last justified current state; no `absent` or `crashed`. No-justified-state case **OPEN**. |
| C-02 | Nonzero exit while stopping projected `absent`. | Replaced with the proposed shared P4-A eligibility and Cases A–D: no `absent` while a failure classification remains unresolved as to establishment, threshold, publication, or public value, or a prior public failure remains unresolved; a bare fatal observation is neither a failure classification nor a bar; S5 non-application is proposed and S5-scope remains OPEN; an established terminal is handled through Case A's CP-1 and active-terminal guards; P4-A-excluded outcomes remain **OPEN** where OD-4 does not uniquely determine them. |
| C-03 | Public `outcome_unknown` for confirmed exits. | Made internal only; confirmed exit/no orderly shutdown uses OD-4; no-exit case preserves the last justified current state, with the no-justified-state case **OPEN**. |
| C-04 | Blanket handshake `unavailable` for missing mandatory capability. | Replaced with OD-6 categories; path classification remains **OPEN** (U8). |
| C-05 | `quarantined` required as a public projection. | Made internal-only; no public state or enum change. |
| C-06 | Diagnostics contradicted O1/O2/O19 state selection. | Diagnostic and projection roles separated; explicit publication rule required. |
| C-07 | Records lacked generation identity. | Generation binding required; exact schema **OPEN**. |
| C-08 | No terminal latch or late-evidence rule. | Immutable first terminal and retention required; non-terminal effects remain **OPEN**. |
| C-09 | “Known action cause” lacked OD-8’s standard. | Replaced by A–E causation plus F assessment; `terminated_by_action` remains internal, requires C-probable or C-established, and retains that qualifier; U5 **OPEN**. |
| C-10 | `seq` had no allocator or event/observation distinction. | Observation order separated from event order; allocator **OPEN**. |
| C-11 | O3 and P4 assigned different public outcomes to the same unexpected exit. | OD-5 classifies `crashed` internally only when exit is observed and orderly shutdown is not established. OD-4 controls public projection through the proposed P4-A and Cases A–D, which separate terminal, fatal-observation, and prior-public-state predicates. No new public `crashed` route; prior public `crashed` preservation is CP-1 OPEN. |
| C-12 | P4 unconditionally projected `absent` on orderly shutdown, even when an established prior terminal or prior fatal failure existed. | Prior terminal handled through Case A's CP-1 and active-terminal guards under OD-3/S5a; `absent` now requires the proposed P4-A predicate. S5 remains retained verbatim and is not accepted policy; its non-application as a public-projection bar is proposed, with S5-scope OPEN; a bare fatal observation does not block OD-4's orderly-shutdown `absent` row. P4-A-excluded branches remain **OPEN** where OD-4 does not uniquely determine the replacement projection. |
| C-13 | P4 described `crashed` as an additional exit fact. | Described as an internal classification; the observed `ExitFact` remains distinct. |
| C-14 | P4, S8, and O4 used inconsistent `absent` conditions and conflated a fatal observation, a justified internal classification, a previously published public state, and an established terminal. | Added proposed P4-A eligibility gates used by S8, O4, P4, O3, and O5; replaced the bare fatal-observation exclusion with a determination-required gate, split P4 into Cases A–D, and split Case C into public-failure, internal-only, and active-state subcategories, with no-prior-failure routed separately to P4-D. Incompatible prior public failures route through ordinary OD-4 handling without being misclassified as active; active `degraded` is confined to C-C. Case A follows OD-3/S5a and adds explicit CP-1 and active-terminal guards without defining the terminal set. Retained S5 is preserved but is not a projection gate; that non-application is proposed, and S5-scope is tracked OPEN. Fatal and non-fatal classifications use the same three-way establishment determination and three-way publication determination, with threshold settlement separated from evidence application; prior public `crashed` preservation is CP-1 OPEN. Genuinely unresolved classification and preservation interactions remain **OPEN** rather than being closed by default. |
| E7–E9 | Not in the original P17 evidence-loss table. | Added and retained as unresolved findings. |

## 13. Status

This is a corrected proposal awaiting independent review. The original P17 remains frozen and unapproved. OD-4's orderly-shutdown outcome is accepted; P4-A and its non-application of retained S5 as a public-projection bar are proposed, not independently accepted policy. S5-scope remains OPEN and is not a gate. CP-1, the prior public `crashed` preservation interaction, remains OPEN.

The correction gives establishment and publication independent three-way determinations and routes terminal values through explicit CP-1 and active-state guards; it does not decide CP-1 or define the terminal set.

Two scope limitations remain explicitly **OPEN**: selection among multiple prior failure classifications, and treatment of a separate non-failure, non-active prior public provider-state value in C-B. This revision neither adopts a multi-classification selection rule nor expands C-B's failure-value routing to every public state.

A qualifying independent review must bind to this exact corrected artifact by path, revision, line count, and SHA-256. The plan-level consistency assessment does not satisfy that requirement. The revision does not amend contract rev3, the decision register, source, tests, or implementation. U1–U8 remain open.

No public provider state was added, removed, or renamed. `outcome_unknown`, `terminated_by_action`, and `quarantined` are not public provider states. No technical mechanism in this draft is accepted policy merely because it is included here.

**P17 REMAINS UNAPPROVED. C1 REMAINS UNAUTHORIZED.**

**No implementation was performed or authorized.**
