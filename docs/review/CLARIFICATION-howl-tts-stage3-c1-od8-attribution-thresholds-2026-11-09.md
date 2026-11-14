# C1 / P17 — OD-8 Clarification: Association vs Attribution (Read-Only)

**Date:** 2026-11-09

**Status:** Clarification only. OD-8 remains **UNRESOLVED**. No OD or U disposition is recorded, and no register change is proposed.

**Basis:** The OD-8 assessment at `docs/review/ASSESSMENT-howl-tts-stage3-c1-od8-causality-2026-11-09.md` (SHA-256 `64447be9…`) and the independent assessment received afterward. This memo does not replace the owner's independent verification of that report.

## 1. Agreed starting point

The existing assessment's central conclusion is accepted as the starting point for this clarification:

> Action record + exit fact + ordering establishes temporal association, not necessarily causation.

A concrete illustration:

1. The process begins exiting because of an independent failure.
2. The supervisor issues `Terminate`.
3. The process reports a nonzero exit.
4. The supervisor observes the exit while its state is `stopping`.

All four observations can be true without the termination action causing the exit. That is why `providerState == stopping` cannot serve as a causal test.

The following points are preserved:

- Action identity and process-instance identity must be independently established.
- Action request, attempted action, successful action, and completed action remain distinct.
- Exit observation must preserve the actual available exit facts, not only a generic error.
- Ordering is necessary for causal reasoning but insufficient by itself.
- Competing explanations must be considered.
- Late evidence may affect the strength of an attribution without rewriting OD-3's immutable terminal decision.
- A successful termination API result proves something about the operation's result; it does not prove that the action caused the exit.

## 2. The clarification in one sentence

**OD-8 should distinguish an association test from an attribution test: the first establishes only that an action and an exit are meaningfully linked in time and instance; the second establishes that the action caused the exit.**

The current alternative A is best read as a necessary-condition floor for association. It is not a sufficient test for causal attribution.

| Test | Question answered | Does it answer OD-8's causal question? |
|---|---|---|
| **Association test** | Was a specific action observed, was an exit observed, were they bound to the same instance, and did the action precede the exit? | No. It establishes A and B only. |
| **Attribution test** | In addition to association, is there evidence that the action operated through a mechanism that produced this exit, with competing explanations addressed? | Yes, when all required evidence is present. |
| **Unresolved** | Is any required element missing, contradictory, delayed, or not bindable to the instance? | No; it must remain unresolved rather than being assigned a cause. |

## 3. Threshold 1 — Association evidence

Association is the minimum needed to say:

> An exit was observed after a specific action was requested or attempted for the same process instance.

It is not a claim that the action caused the exit.

### 3.1 Required elements

| Element | Minimum content | Why it is needed |
|---|---|---|
| **A1. Action observation** | A specific action request or attempt is recorded, with its kind (terminate, kill, shutdown, cleanup, cancellation) and its target instance. | "The provider was stopping" is a state, not an action. A generic state change cannot identify which action is being associated with the exit. |
| **A2. Exit observation** | An actual process exit is recorded, with the available exit status/error and observation context. | A request, a successful API return, or a state change is not an exit fact. The exit must be observed as an event. |
| **A3. Instance binding** | The action and exit are bound to the same authoritative process instance. A generation/incarnation token is the target concept under OD-2; a PID may corroborate but cannot alone survive replacement or reuse. | Without this, the association may be between an action on one process and the exit of another. |
| **A4. Ordering** | The action observation precedes the exit observation in a recorded order for that instance. A monotonic sequence is preferable to wall-clock time. | Without ordering, even temporal association is conjecture. |

### 3.2 What association does not require

Association does not require proof of a causal mechanism. That is why it is a separate, lower threshold. It does, however, require enough identity and ordering that the word "association" is not doing hidden causal work.

### 3.3 What association cannot establish

Even when A1–A4 are all present:

- the action may have been requested after the process had already begun exiting;
- the action may have been rejected or timed out;
- the exit may have been caused by an independent failure;
- the exit may have been a separate shutdown-time failure;
- a fatal observation may explain the session while the action explains only the later exit, or vice versa.

Association therefore remains **B**, not **C**.

### 3.4 Current-source status

The current source does not fully satisfy A1–A4:

- no durable action record identifies the action kind and target;
- `Kill` is an alias of `Terminate`;
- `waitErr` is an error, not a structured exit fact;
- `stopping` is a mutable state check, not an ordering record;
- no supervisor generation is retained with the action or exit.

The source can at most produce a state-based association in some interleavings. It cannot produce a fully recorded association on demand.

## 4. Threshold 2 — Causal-attribution evidence

Attribution is the additional evidence needed to say:

> This process instance exited because of this deliberate action.

It requires the association elements **plus** the following.

### 4.1 Additional required elements

| Element | Minimum content | Why it is needed |
|---|---|---|
| **C1. Association satisfied** | A1–A4 are all present. | Causation presupposes that the action and exit are actually linked to the same instance and ordered. |
| **C2. Action outcome** | The record distinguishes requested, accepted, attempted, completed, rejected, and timed out. A successful API return is not automatically acceptance or completion. | An action request is not an action result. A rejected or timed-out action cannot be treated as a completed cause. A `nil` return may mean the done channel was observed, not that the action produced the exit. |
| **C3. Exit fact** | The exit observation preserves the available exit code, signal, error, and observation context. A generic error is not a substitute for an exit status. | The exit fact is the effect to be attributed. Without its actual content, the attribution cannot distinguish termination modes or check consistency with the claimed mechanism. |
| **C4. Mechanism evidence** | Evidence that the action operated through a mechanism capable of producing this exit, and that the mechanism operated in this case. Examples: an OS-reported kill/signal outcome; a record that the action reached the target process; a supervisor observation that the targeted process ended after the action's signal. For a fatal observation, evidence of how the fatal condition produced the exit. | This is the step from B to C. Ordering and an action record show sequence, not production. |
| **C5. Competing-explanation disposition** | A pre-existing exit, independent failure, fatal observation, shutdown-time failure, or second action is recorded and either excluded by evidence or left explicitly unresolved. | A causal claim is not established by ignoring an explanation that fits the same observations. |
| **C6. Uncertainty statement** | The record states whether the attribution is established, probable, or unresolved, and on what evidence. | The standard must not turn missing evidence into certainty. A probabilistic mechanism can support a probabilistic attribution, but the record must say so. |

### 4.2 What attribution does not require

Attribution does not necessarily require a counterfactual experiment or a perfect OS-level cause on every platform. It does require enough mechanism and exclusion evidence to distinguish the action-caused explanation from the alternatives that the available facts permit.

A stricter standard is compatible with this one. A standard that stops at A1–A4 is not.

### 4.3 Fatal observations

For a fatal observation, replace "action identity" with the identity of the fatal observation and the policy that responded to it. The same requirements apply:

- same instance;
- identifiable observation and response;
- outcome and exit fact;
- ordering;
- mechanism from the fatal condition to the exit;
- competing explanations;
- explicit uncertainty.

A fatal observation followed by an exit is not, by itself, proof that the observation caused the exit.

## 5. Contradictory, incomplete, and delayed evidence

The standard must define what happens when the evidence does not fit a clean story.

### 5.1 Decision table

| Evidence pattern | Association? | Attribution? | Result |
|---|---|---|---|
| Action request, exit observed, same instance, action first, but no mechanism or competing-explanation record. | Yes, if A1–A4 are recorded. | No. | **B only**; causal result unresolved. |
| Action request and exit observed, but no same-instance binding. | No. | No. | **Unresolved**; do not associate across instances. |
| Action request and exit observed, but no ordering. | No complete association. | No. | **Unresolved**; temporal relation is conjecture. |
| Exit observed before the action. | No action-caused association. | No. | The later action cannot have caused that earlier exit; record as a separate event if relevant. |
| Action rejected or timed out, then exit. | At most a temporal association with the request. | No, unless a separate attempted action and mechanism are shown. | **Unresolved**; a rejected action is not a completed cause. |
| Action completed, exit status present, mechanism consistent, no competing explanation. | Yes. | Yes, subject to C6. | **C** for that instance. |
| Action completed, exit status present, but a pre-existing fatal observation is not excluded. | Yes. | Not established. | **Unresolved**; retain both explanations. |
| Two actions and one exit, with no identity/order/mechanism separating them. | Possibly. | No unique attribution. | **Unresolved**; do not choose the first or last action by convenience. |
| Late evidence contradicts the mechanism. | May remain. | Attribution weakens. | **Unresolved or downgraded**; the immutable terminal decision is not rewritten, and the late evidence is retained. |

### 5.2 Rules

1. **Do not resolve by class order.** A fatal observation does not automatically outrank an exit fact, and an exit fact does not automatically erase a fatal observation.
2. **Do not resolve by first-seen failure.** The first stored failure is not necessarily the first event.
3. **Do not resolve by last-writer state.** A later state write is not a causal finding.
4. **Do not silently drop evidence.** Contradictory or delayed evidence must remain recordable, even after a terminal decision.
5. **Do not turn ambiguity into a projection.** `crashed`, `absent`, and `outcome_unknown` are not interchangeable. The first two are provider-state or disposition projections; the last is an internal acknowledgment that the evidence does not establish a terminal outcome. Its public projection remains undecided.

## 6. When attribution must remain unresolved

Attribution must remain unresolved when any of the following is true:

- A1–A4 are not all present;
- the action and exit cannot be bound to the same instance;
- the action outcome is unknown, rejected, or timed out without a separate completed attempt;
- the exit fact is only a generic error and cannot support the claimed mechanism;
- no mechanism evidence connects the action to this exit;
- a competing explanation cannot be excluded;
- the evidence is contradictory, incomplete, or delayed in a way that leaves more than one explanation standing;
- late evidence weakens the mechanism below the attribution threshold.

In those cases:

- do **not** conclude that the action caused the exit;
- do **not** automatically conclude `crashed` or `absent`;
- record the result as unresolved/`outcome_unknown` internally where that distinction is available;
- leave the public projection to the separate OD-4b/OD-5 decisions;
- preserve the evidence for later review.

This is the point of the clarification: **unresolved is a legitimate result, not a failure to choose a winner.**

## 7. U5 remains open

U5 is the application case: nonzero exit after a deliberate kill. The clarification does not close it.

| Required element | Current source | Consequence for U5 |
|---|---|---|
| Same instance | Process pointer and PID exist; no supervisor generation is retained with the action/exit. | Partial; replacement and reuse cannot be excluded from evidence alone. |
| Action identity | `Kill` aliases `Terminate`; no action ID or kind is stored. | Absent. |
| Action outcome | `Terminate` returns nil/error/timeout; several callers discard it. | Absent as structured evidence. |
| Exit fact | `waitErr` is stored; no exit code, signal, or observation timestamp. | Partial only. |
| Ordering | No timestamp or sequence; the only signal is `providerState == stopping` at observation time. | Inferred, not recorded. |
| Mechanism | No signal/cause record or action-to-exit link. | Absent. |
| Competing explanations | A pre-existing exit, independent failure, or shutdown-time failure is not separately retained. | Cannot be excluded. |

**Conclusion:** U5 cannot presently receive a unique, source-justified causal classification. It may have an action observation and an exit observation, but not the evidence needed for C. U5 remains **OPEN**, and this memo does not decide its public projection.

The four-step illustration from the independent assessment is exactly the U5 hazard: all four observations can be true while the termination action is not the cause.

## 8. Proposed non-binding wording for OD-8 clarification

This is not a register update and not a disposition. It is a wording option for the owner to consider if OD-8 is later revised.

> **OD-8 clarification (non-binding):**
>
> 1. **Association, not causation.** An action and an exit are associated only when a specific action request/attempt, an actual exit observation, the same process-instance identity, and a recorded action-before-exit ordering are all present. Association never establishes causation.
> 2. **Causal attribution.** Attribution requires association plus a recorded action outcome (requested, accepted, attempted, completed, rejected, or timed out), the actual available exit facts, evidence of a mechanism by which the action produced this exit, and a disposition of competing explanations. A successful API return is not proof of causation.
> 3. **Unresolved.** If any attribution element is missing, contradictory, delayed, or not bindable to the instance, the result remains unresolved/`outcome_unknown`; it is not converted into `crashed`, `absent`, or an action-caused termination by assumption.
> 4. **Terminal and late evidence.** The immutable first terminal decision under OD-3 is not rewritten by late evidence. Late evidence remains recordable and may strengthen, weaken, or leave the attribution unresolved. Publisher authority, evidence ownership, synchronization, and projection effects remain subject to OD-1/OD-3's unresolved exclusions.
> 5. **Fatal observations.** The same association and attribution requirements apply to a fatal observation and the response to it; a fatal observation followed by an exit is not automatically the cause of that exit.

This wording keeps alternative A's conservatism while separating the association floor from the attribution test. It does not select an owner disposition.

## 9. What this clarification does and does not change

### It does

- separate the minimum association evidence from the additional attribution evidence;
- define the treatment of contradictory, incomplete, and delayed evidence;
- define when attribution must remain unresolved;
- preserve the independence of U5 and every other U case;
- preserve OD-3's immutable terminal decision while allowing the evidence record to grow;
- keep the thirteen provider states frozen and P17 unapproved.

### It does not

- accept, reject, or defer OD-8;
- change OD-1 (deferred), OD-2/OD-3 (accepted), OD-4–OD-9 (unresolved), or U1–U8 (open);
- decide that U5 or any other case is action-caused;
- decide `crashed` semantics, public projection, quarantine, or compatibility taxonomy;
- decide evidence-record ownership or publisher authority;
- modify the decision register, P17, the corrected r3 response, the lifecycle decision, or contract rev3;
- create rev4, P18, a contract amendment, or an implementation;
- run tests, builds, benchmarks, or runtime experiments;
- claim that the current source satisfies the proposed standard;
- substitute for the owner's independent verification of the assessment report.

## 10. Verification note

The existing OD-8 assessment is unchanged at:

`64447be9d166c3d7d9b6c09ce332180d5271168513f9f24b6aa91284472cb17e`

The decision register is unchanged at:

`4a2a24b908cbb382f05711d10750bc156e0abf7358eec1b7e98ef50116663a16` (244 lines)

The frozen artifact hashes remain:

- P17: `01d887c1…`
- Corrected r3 response: `60ca1cb5…`
- Lifecycle decision: `aa136bc9…`
- Contract rev3: `ed0da1a7…`

The owner's independent inspection and hash recomputation remain outstanding. This memo does not claim to supply that independent verification.

**OD-8 remains UNRESOLVED. No register change or implementation is warranted.**
