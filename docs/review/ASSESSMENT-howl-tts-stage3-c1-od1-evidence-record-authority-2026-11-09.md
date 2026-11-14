# C1 / P17 — OD-1 Evidence-Record Authority Assessment (Read-Only)

**Date:** 2026-11-09

**Status:** Assessment only. **OD-1 remains DEFERRED.** No OD-1 disposition is recorded here. P17 remains unapproved; U1–U8 remain OPEN; contract rev3 remains frozen; C1 implementation remains unauthorized.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `f9f0a1710776f2ca8870c87d9cfd4f6f6d4c61fc9b64f4509db9bf405e389960`, 409 lines.
- Dispositions: **8 ACCEPTED** (OD-2 through OD-9); **1 DEFERRED** (OD-1); **0 UNRESOLVED**.
- P17 is a frozen author proposal, not approved; contract rev3 is a proposed revision, not accepted.

This assessment addresses OD-1 only. It does not decide OD-1, revise P17, amend contract rev3, update the register, change source, run tests, or recommend implementation as the next step.

## 1. Exact OD-1 wording

The register asks:

> **OD-1 — Evidence-record authority**
>
> **Question:** Who owns the authoritative evidence record, and is it diagnostic-only or projection-supporting?
>
> **Status:** DEFERRED.
>
> **Proposed default:** None proposed; deferred without an authority decision.
>
> **Downstream:** R09, R10, U1–U7.

The exact wording is quoted for decision precision. It is not a new disposition.

## 2. What OD-1 actually asks

OD-1 is two related governance questions, not a request to design a logging subsystem:

1. **Ownership and authority:** which component or role is accountable for the evidence record, and what may that record authoritatively establish?
2. **Purpose:** is the record diagnostic-only, or may it support provider-state projection?

The second question is the more consequential one. A record can exist, be append-only, and still be diagnostic-only. “Append-only” does not make a record authoritative, and “authoritative” does not mean every entry is complete, ordered, or causally sufficient.

OD-1 does **not** ask:

- which file format or storage engine to use;
- which source functions must change;
- whether to add a provider state;
- whether to approve P17;
- whether to implement C1.

Those are downstream questions. OD-1 decides what the evidence record is *for* and who may rely on it.

## 3. Evidence classification

### 3.1 Observed

Directly read from the current governance record:

- OD-1 is DEFERRED and has no recorded rationale or selected authority model.
- P17 is unapproved; U1–U8 are open; contract rev3 is not accepted.
- OD-9 is accepted as Option A: P17 is not approved while OD-1 remains deferred and U1–U8 remain open without determined outcomes or explicit unknown-case rules.
- P17 §1 requires all observations to be preserved in an append-only log, and §7 applies a precedence table to that preserved log.
- P17 §8 states that the evidence-loss paths must be fixed before precedence can be trusted.

### 3.2 Source-derived

The following are established by the existing source assessment, not by runtime experiments in this review:

- There is no append-only evidence record, sequence allocator, or stable process-instance identity stored with failures, actions, and exits.
- `Failure` carries code, message, retryability, field, and detail only; it has no instance, sequence, timestamp, phase, severity, action, or causal field.
- The supervisor retains one `lastFailure` and one `waitErr`; `SetState` overwrites `lastFailure`.
- `failAll` preserves only the first read failure and delivers that same failure to all pending requests.
- `Terminate` records no intent, start, completion, or result, and its result is ignored at several call sites.
- `waitForReader` returns nothing on join or timeout, so the caller cannot distinguish those outcomes.
- There is no terminal latch, terminal revision record, quarantine record, or ownership-uncertainty field.

These facts matter because they describe the gap between P17's model and the current evidence surface. They do not decide whether a new record should be authoritative.

### 3.3 Reported

- The owner accepted OD-2 through OD-9; those are recorded governance directions.
- The prior assessment reports R09–R13 and U1–U8 as open or narrowed.
- P17 itself claims an append-only evidence model but does not state who owns it or whether it supports projection.

### 3.4 Inferred

- If the record is diagnostic-only, P17's §7 cannot claim authority to resolve provider-state projections from it.
- If the record is projection-supporting, the owner must accept a stronger set of obligations: identity, sequence, action/exit linkage, retention, publication authority, and conflict resolution.
- A decision to defer OD-1 is not neutral: it leaves P17's current projection claim unsupported.

### 3.5 Unestablished

- Whether the record is intended for diagnosis, projection, or both.
- Who owns the record and who may publish from it.
- Whether the record must be durable, complete, or ordered.
- Whether “absence of an entry” means “no event,” “not observed,” or “lost.”
- Whether a projection-supporting record must be retained across process lifetimes.
- Whether OD-1 requires any contract-level rule beyond the frozen provider-state set.

## 4. What P17 currently assumes

P17 does not merely propose a log; it uses the log as the basis for projection:

- §1 constraint 1: all observations are preserved in an append-only log.
- §2: `Observation`, `Action`, and `EvidenceLog` are defined as records.
- §7: precedence selects a projection from the preserved log.
- §7's table maps evidence combinations to `absent`, `crashed`, `incompatible`, `unavailable`, `outcome_unknown`, and internal `terminated_by_action`.
- §8: evidence-loss findings E1–E6 must be fixed before precedence can be trusted.
- §10: the reviewer is asked to confirm five projections, including `quarantined`, are reachable.

That is a projection-supporting posture in substance, but P17 never states the authority. OD-1 is therefore logically prior to P17's §7 claim. A record that is only diagnostic cannot carry the weight P17 places on it.

## 5. Authority models

The following models are governance options, not recommendations. They differ in what the record may establish, not merely in where it is stored.

### 5.1 No authoritative record — status quo

**Description:** Current in-memory state and ad hoc fields remain the only evidence surface. There is no append-only record, no declared owner, and no authority to reconstruct a history.

**What it establishes:** Individual facts currently retained by the supervisor or client, subject to overwrite and loss. It does not establish an ordered history.

**Consequences:**

- P17 §7 has no authoritative record to consume; its precedence table is unsupported.
- R09 and R10 remain open; U1–U7 cannot receive a defensible evidence-based outcome.
- P17 cannot be approved as written.
- No contract amendment or implementation follows.
- The status quo is not neutral: it preserves the existing evidence losses and leaves the projection claim ungrounded.

### 5.2 Diagnostic-only record

**Description:** A record is maintained for debugging, audit, or operator diagnosis, but it is not the authority for provider-state projection. It may be append-only, best-effort, or bounded; its completeness is not a contract guarantee.

**What it establishes:** Diagnostic context. It does not establish the public provider state, terminal outcome, or causal attribution.

**Consequences:**

- P17 §7 must be narrowed or withdrawn as a projection authority. The selection table may remain a design note, but it cannot claim to resolve state from the log.
- R09 and R10 are not closed by the record alone; U1–U7 still need other evidence or explicit unknown dispositions.
- The frozen provider-state set remains unchanged.
- No contract amendment is required merely to keep a diagnostic log.
- A future implementation would still be a source change and remains unauthorized now.
- Diagnostic-only is compatible with OD-2 identity and OD-3 terminal rules, but it does not satisfy them by itself: a diagnostic record may still need identity and ordering to be useful.

**What it cannot do:** It cannot turn an unobserved exit into a confirmed crash, an action into a cause, or an absent entry into proof that no event occurred.

### 5.3 Projection-supporting record

**Description:** The record is authoritative for the facts from which a provider-state projection may be selected. It must preserve the observations, actions, and exit facts relevant to the projection, with enough identity, order, and linkage to support the declared rules.

This is not merely “a better diagnostic log.” It is an authority claim. It implies obligations:

- a single accountable owner for the record;
- a defined instance identity and sequence domain;
- action intent, start, completion, and result where action causation matters;
- exit facts separated from exit causes;
- a retention and loss policy;
- a terminal-publication and late-evidence rule consistent with OD-3;
- a rule for conflicts and missing entries;
- a clear boundary between internal evidence and public provider states.

**Consequences:**

- P17 can claim projection authority only after its model and the record semantics are reconciled; the current text is not sufficient.
- R09 and R10 can be addressed by the record design, but neither is automatically closed. A complete record still needs a causal standard (OD-8) and a defined publisher/linearization rule (OD-3).
- U1–U7 may become reviewable against preserved facts, but runtime interleavings still need actual observation. A record can preserve an execution; it cannot prove which execution occurred without one.
- No new public provider state is implied; `quarantined` remains internal per OD-7.
- Contract rev3 may need alignment with evidence/publication rules, but OD-1 does not itself amend it. The frozen thirteen-value provider-state set remains unchanged.
- Implementation would be substantial and is not authorized by this assessment.

### 5.4 Split record — diagnostic log plus projection record

**Description:** A diagnostic log and a separate authoritative projection record coexist. The diagnostic log is best-effort and broad; the projection record is narrow, complete for its declared scope, and authoritative.

**Consequences:**

- The authority boundary must be explicit: which record may be cited for a projection, and what happens when they disagree.
- The split can reduce the temptation to make every diagnostic message projection-authoritative.
- It adds duplication, retention, and consistency obligations; two records can diverge.
- It does not avoid OD-3's publisher/linearization question or OD-8's causation question.
- It is a variant, not a way to defer the authority choice.

## 6. Ownership variants

OD-1 also asks **who owns** the record. Ownership is not the same as physical storage. At minimum it must identify who may append, who may publish, who resolves conflicting entries, and who is accountable for retention.

| Owner model | Strength | Principal risk | Authority consequence |
|---|---|---|---|
| **Supervisor-owned** | Natural owner of process lifecycle, termination actions, and exit facts. | Client-side handshake/protocol observations must be submitted or merged; gaps at the boundary. | Can be authoritative for process/exit facts, but not automatically for handshake or request context. |
| **Client/session-owned** | Natural owner of handshake, protocol, request, and transport context. | Process exit and termination facts live in the supervisor; client lifetime is not process lifetime. | Can be authoritative for session facts, but not automatically for process outcomes. |
| **Shared record with one linearizer/publisher** | Can cover client and supervisor facts under one authority. | Requires a defined single writer or linearization point; otherwise “shared” means no authority. | Potentially authoritative across both domains, but only with explicit ordering and publication rules. |
| **No single owner** | Low coordination cost. | No accountable resolver; conflicting entries cannot be adjudicated. | Not authoritative; effectively diagnostic-only or no record. |
| **Separate diagnostic collector** | Useful for audit and debugging. | Not part of the lifecycle authority. | Diagnostic only unless separately granted projection authority. |

A shared record without a single accountable publisher or linearization authority is not a projection-supporting record merely because it is append-only.

## 7. Consequence matrix

| Model | Record authority | P17 effect | U1–U8 effect | Contract effect | Implementation effect |
|---|---|---|---|---|---|
| **No authoritative record** | None | §7 unsupported; P17 unapproved | U1–U7 remain unresolved | No amendment | None authorized |
| **Diagnostic-only** | Diagnostic context only | §7 must be narrowed or withdrawn; P17 unapproved as written | U1–U7 still need other evidence or explicit unknown | No new provider state | Any future log is a source change; not authorized |
| **Projection-supporting** | Authoritative within declared scope | P17 must be reconciled; can claim projection only after design | Evidence may support outcomes, but runtime interleaving still required | No new public state; possible rule alignment, no automatic amendment | Substantial future work; not authorized |
| **Split record** | Authoritative only for the projection record | P17 must name the authoritative record | Same as projection-supporting, plus consistency questions | No new public state | More complex future work; not authorized |

The matrix is not a scorecard. Diagnostic-only is not “safe” if P17 continues to claim projection authority; projection-supporting is not “approved” merely because the owner selects it.

## 8. What a defensible OD-1 decision must state

A decision that leaves any of these unanswered is not yet a usable authority model:

1. **Purpose:** diagnostic-only, projection-supporting, or split.
2. **Owner:** the accountable component or role.
3. **Authoritative scope:** which facts the record may establish, and which it may not.
4. **Instance identity:** how entries bind to the process generation (OD-2).
5. **Ordering:** who allocates sequence and how concurrent entries are linearized.
6. **Action linkage:** what action record is required before an exit may be attributed to an action (OD-8).
7. **Terminal rule:** whether a published terminal is immutable and how late evidence is retained (OD-3).
8. **Retention and loss:** what must be retained, for how long, and what a missing entry means.
9. **Publication boundary:** who may project a provider state from the record.
10. **Contract boundary:** whether any provider-state or contract rule changes; no new public state is assumed.
11. **P17 consequence:** whether P17 §7 may claim projection authority or must be narrowed.
12. **Implementation boundary:** that selecting a model does not authorize source changes, tests, builds, or C1.

A record that cannot answer “what does absence mean?” and “who resolves a conflict?” is not authoritative in the sense P17 needs.

## 9. Decision readiness

OD-1 can be assessed and decided independently of U1–U8 because it is an authority question, not an outcome question. The owner does not need to know whether U1 is `absent` or `crashed` to decide whether the record may support projection at all.

It cannot be decided by source inspection alone. The source establishes that the current record is absent and lossy; it cannot choose whether the owner wants a diagnostic record or an authoritative one.

The dependency direction is:

1. **OD-1** — authority and purpose;
2. **R09/R10** — record, ordering, and causal rules;
3. **OD-3/OD-5/OD-6/OD-7/OD-8** — terminal, crash, compatibility, quarantine, and causation rules;
4. **U1–U8** — individual outcomes or explicit unknown rules;
5. **P17 approval** — separate decision, still gated by OD-9.

OD-1 is upstream of R09, R10, and U1–U7. U8 is primarily a compatibility-taxonomy question and does not depend on OD-1 for its classification, though P17 approval still depends on both.

Deferring OD-1 is possible, but it has a consequence: P17's §7 remains unsupported, U1–U7 cannot be finalized, and OD-9's gate remains unsatisfied. Deferral is not a neutral “wait and see” if P17 approval is being considered.

## 10. Ambiguities

- **“Authoritative”** is undefined. It could mean the sole source of truth, the best available evidence, or merely a durable record. OD-1 needs one meaning.
- **“Diagnostic-only”** is undefined. A diagnostic log may be best-effort or append-only; neither makes it projection-authoritative.
- **“Who owns”** could mean storage owner, writer, publisher, or adjudicator. These are different authorities and can belong to different components.
- **Scope:** “projection-supporting” does not say whether every fact, only selected facts, or only terminal facts are authoritative.
- **Absence semantics:** no document says whether a missing entry means no event, an unobserved event, or a lost event.
- **Retention:** no period, durability, or loss policy is stated.
- **Contract relationship:** it is unclear whether a projection-supporting record requires contract text beyond the frozen provider-state set. This assessment does not assume an amendment.
- **P17 relationship:** P17 assumes a projection-supporting record but does not state the authority. A decision to make the record diagnostic-only therefore requires a P17 text change, not merely a new log.

## 11. Recommendation

**No authority model is selected by this assessment.** The owner must choose between diagnostic-only, projection-supporting, and a split model; OD-1 cannot be decided by inference from P17.

The decision should be made in this order:

1. **Decide whether P17 §7 is intended to be authoritative.** If it is, the record must be projection-supporting. If it is not, the record is diagnostic-only and P17's §7 claim must be narrowed.
2. **If projection-supporting, accept the obligations in §8 before selecting an owner.** A record without identity, ordering, action linkage, retention, and publication authority cannot satisfy the claim.
3. **If diagnostic-only, state explicitly that U1–U7 remain unresolved unless separate evidence establishes them.** Diagnostic preservation is not an outcome rule.
4. **If split, define which record is authoritative and how conflicts between the two are resolved.**
5. **Do not treat any choice as P17 approval.** OD-1 selects the authority model; it does not approve P17, amend contract rev3, or authorize implementation.

The one unconditional conclusion is this: **P17's §7 cannot be treated as authoritative while OD-1 is deferred, and a diagnostic-only record cannot support it as written.** That is not a vote for projection-supporting; it is the consequence of P17's own claim.

A defensible OD-1 decision should also state that it does not:

- approve P17;
- amend contract rev3;
- add a public provider state;
- close U1–U8;
- authorize source changes, tests, builds, benchmarks, or runtime experiments;
- authorize C1 implementation.

## 12. Limitations and verification

This is a read-only assessment. No existing artifact was modified. The decision register remains at `f9f0a171…`, 409 lines; P17 remains `01d887c1…`; contract rev3 remains `ed0da1a7…`; no source, test, configuration, or generated file was changed.

**Limitations:** this assessment relies on the existing P17 assessment and register for source-derived facts; it did not re-run the full source audit. It does not execute any U case, observe any interleaving, or verify platform-specific behavior. It does not decide OD-1, and it does not establish that any implementation conforms to P17 or the accepted decisions.

**OD-1 remains DEFERRED. P17 remains UNAPPROVED.**
