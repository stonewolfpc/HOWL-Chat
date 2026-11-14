# C1 / P17 — OD-8 Requirements Challenge: Ordering and Causal Mechanism (Read-Only)

**Date:** 2026-11-09

**Status:** Review only. OD-8 remains **UNRESOLVED**. No OD or U disposition is recorded, and no register change is proposed.

**Inputs:**

- Clarification memo: `docs/review/CLARIFICATION-howl-tts-stage3-c1-od8-attribution-thresholds-2026-11-09.md`, SHA-256 `6dcc3b35…`.
- OD-8 assessment: `docs/review/ASSESSMENT-howl-tts-stage3-c1-od8-causality-2026-11-09.md`, SHA-256 `64447be9…`.
- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `4a2a24b9…`, 244 lines.

## 1. Review result

The two-threshold separation is sound and should be retained:

- **Association** answers whether a specific action and an observed exit belong to the same process instance and are ordered.
- **Causal attribution** additionally requires evidence of a mechanism connecting the action to the exit and a disposition of competing explanations.

However, both challenged requirements need refinement before they can support an owner decision:

1. **A4 (ordering)** currently requires that the action *observation* precede the exit *observation* in a recorded order. That is not necessarily the same as the action *event* preceding the exit *event*. Observation order, log order, and wall-clock timestamps can all misrepresent event order.
2. **C4 (mechanism)** currently permits evidence that the action operated through a mechanism "capable of producing" the exit. That wording can be satisfied by a plausible narrative. The standard needs evidence that the mechanism actually operated on this instance, plus a stated way the claim could fail.

The correct response is to refine these requirements, not to abandon them. Ordering and mechanism remain necessary for causal attribution; the clarification is about **what quality of evidence satisfies them**.

The safeguards from the prior review remain in force:

- a successful action result is not causal proof;
- a fatal observation is not automatically the cause of termination;
- late evidence can change attribution confidence without rewriting OD-3's immutable terminal decision;
- U5 remains open;
- no public projection or implementation is decided.

## 2. Challenge to A4 — ordering

### 2.1 The current requirement

The clarification memo's A4 says:

> The action observation precedes the exit observation in a recorded order for that instance. A monotonic sequence is preferable to wall-clock time.

This is directionally correct but under-specified in four ways.

### 2.2 Observation order is not event order

The requirement orders **observations**, but the causal question concerns **events**. These are different:

| Time concept | Meaning |
|---|---|
| Action occurrence | When the action actually took effect, or when its operative signal was delivered. |
| Action record | When the action was written to the evidence record. |
| Action observation | When a component noticed the action or its effect. |
| Exit occurrence | When the process actually exited. |
| Exit record | When the exit was written to the evidence record. |
| Exit observation | When a component noticed the exit. |

A recorded action-before-exit order can be produced even when the exit event occurred first:

- the action request is logged promptly, but the process had already begun exiting;
- the exit occurs asynchronously, but its wait result is observed or recorded later;
- an action callback is delayed and written after the exit record;
- a buffered writer flushes records out of event order;
- the action request is recorded before the kill signal is delivered, while the exit happens before delivery.

The current A4 can therefore establish **observation-order association**, not necessarily **event-order association**. Those must be distinguished.

### 2.3 Wall-clock timestamps are not a reliable ordering source

Wall-clock timestamps alone are insufficient because of:

- clock skew between components;
- clock adjustments, including backward jumps;
- timestamp resolution and truncation;
- different clocks for action, exit, and observation;
- asynchronous callbacks and scheduling delays;
- buffering and delayed persistence;
- records written after the event they describe.

A timestamp can corroborate order. It cannot, by itself, establish a total order for a causal claim.

### 2.4 A sequence assigned at record time has the same defect

A monotonic sequence is stronger than a timestamp only if it is assigned at the point that actually orders the events. A sequence assigned when a record is written inherits the observation/record delay. It orders **writes**, not necessarily **events**.

A reliable event-order source needs one or more of:

- a per-instance monotonic sequence assigned by a single authoritative recorder at the point of event occurrence or acceptance;
- an OS/process-level ordering fact that directly establishes the action effect before the exit;
- a causal handshake in which the action's operative message is acknowledged before the exit is recorded;
- a partial order built from happens-before relationships, with unresolved concurrency represented explicitly.

These are possible mechanisms, not implementation requirements. The evidence requirement is: **the record must establish event order, not merely write order.**

### 2.5 Concurrency and overlap

Action and exit can overlap. The process may already be exiting when the termination action is requested. In that case:

- the action request may be observed before the exit is observed;
- the action effect may occur after the exit event;
- neither event may be reliably ordered relative to the other.

A total order may be unavailable. The standard should not manufacture one. It should record a partial order and mark the relation unresolved when the events are concurrent or their order is unknown.

### 2.6 Which action event is the anchor?

A4 says "the action observation" without defining which action event matters:

- the request;
- the attempt;
- the acceptance;
- the operative effect (for example, signal delivery);
- completion.

For **association**, the request or attempt may be enough to identify a temporal relationship, but it is not enough to establish that the action had any effect before the exit. For **attribution**, the relevant anchor is the operative effect, not the request. A request can precede an exit by a long interval while the effect occurs after the exit.

The standard should therefore distinguish:

- **A4-request:** the action request/attempt is ordered before the exit observation;
- **A4-effect:** the action's operative effect is ordered before the exit event.

A4-request can support observation-order association only. A4-effect is what causal attribution needs.

### 2.7 Proposed refinement for A4

A non-binding wording option:

> **A4. Ordering.** The record establishes that the relevant action event occurred before the exit event for the same instance. For association, the request/attempt event may be used, but the result must be labeled observation-order association if event order is not established. For attribution, the action's operative effect must be ordered before the exit. A per-instance monotonic sequence assigned by a single authoritative recorder at the point of occurrence is preferred. Wall-clock timestamps, log order, or a sequence assigned only at record time are corroboration, not proof of event order. If the events are concurrent or their order is unknown, record the relation as unresolved.

This preserves A4's purpose while preventing a recorded order from masquerading as a causal order.

## 3. Challenge to C4 — causal mechanism

### 3.1 The current requirement

The clarification memo's C4 says:

> Evidence that the action operated through a mechanism capable of producing this exit, and that the mechanism operated in this case. Examples: an OS-reported kill/signal outcome; a record that the action reached the target process; a supervisor observation that the targeted process ended after the action's signal.

The second clause is the right idea. The first clause and the examples are too easy to satisfy with a plausible story.

### 3.2 Capability is not operation

A mechanism can be capable of producing an exit without operating in the particular case:

- a kill signal can be sent but not delivered;
- a kill can be delivered after the process has already exited;
- an OS-reported exit status can be consistent with a kill without proving the kill caused the exit;
- a helper-side record can be written before the exit but not establish that the signal produced it;
- a supervisor can observe an exit after a signal even though the exit was already in progress.

C4 must require evidence that the mechanism **operated on this instance**, not merely that it could have.

### 3.3 "Mechanism" can become an unfalsifiable narrative

If any exit can be explained by "shutdown," "termination," or "the process was stopping," then the mechanism requirement adds no discriminating power. A mechanism claim needs:

- a specific pathway;
- observable implications if the pathway operated;
- at least one defeater that would show it did not operate;
- a record of whether that defeater is present, absent, or unresolved.

Examples of defeaters:

- the exit occurred before the operative signal;
- the action was rejected or timed out before effect;
- the action targeted a different instance;
- the exit status is inconsistent with the claimed mechanism;
- an independent fatal condition had already occurred;
- a second action or shutdown path better explains the exit;
- no evidence shows the action reached the process.

A mechanism claim that cannot name a defeater is not an evidentiary claim; it is a narrative.

### 3.4 Exit status is consistent with, not proof of, causation

An OS-reported kill status can be strong evidence, but it still needs context:

- the process may have been dying when the signal arrived;
- the status may reflect the last signal rather than the initiating cause;
- a process can exit nonzero for reasons unrelated to the action;
- a zero exit can follow a fatal observation or a shutdown action.

C4 should require consistency **plus** operation evidence and competing-explanation handling, not status alone.

### 3.5 Direct and indirect fatal chains

A fatal observation can relate to an exit in at least three different ways:

1. **Direct:** the fatal condition itself causes the process to exit.
2. **Indirect through an action:** the fatal condition triggers a termination action, and that action causes the exit.
3. **Independent:** the fatal condition and the exit have separate causes, or the action and the fatal condition both contribute.

The clarification memo's C4 says "evidence of how the fatal condition produced the exit." That is appropriate for the direct case but incomplete for the indirect case. The standard should require the record to distinguish:

- fatal condition → exit;
- fatal condition → action → exit;
- action → exit independent of the fatal condition.

Timing alone does not distinguish these. A fatal observation followed by a kill followed by an exit does not establish which link caused the exit.

### 3.6 Probabilistic mechanisms

Some mechanisms may be probabilistic or platform-dependent. The standard should not force false certainty. It should distinguish:

- **C-established:** mechanism operation is supported by direct evidence and competing explanations are addressed;
- **C-probable:** mechanism operation is supported but not conclusive; the confidence and the missing evidence are stated;
- **C-unresolved:** mechanism operation is absent, contradicted, or indistinguishable from alternatives.

A probabilistic mechanism can support a probabilistic attribution, but the record must say so. It must not be promoted to C-established by confidence alone.

### 3.7 Proposed refinement for C4

A non-binding wording option:

> **C4. Mechanism evidence.** The record must show (a) a specific mechanism capable of producing this exit, (b) evidence that the mechanism operated on this instance, and (c) that the observed exit is consistent with that mechanism. The record must identify at least one defeater that would show the mechanism did not operate, and state whether that defeater is present, absent, or unresolved. A successful action result, a plausible narrative, or a post-hoc correlation is not sufficient. For a fatal observation, distinguish a direct fatal-condition cause from a fatal-condition-triggered action and from an independent exit. If mechanism operation is probabilistic, record the attribution as probable, not established, and state the missing evidence.

This keeps C4's causal role while making it falsifiable and less susceptible to narrative completion.

## 4. Interaction between ordering and mechanism

Ordering and mechanism are related but not interchangeable:

- mechanism evidence can help establish event order, for example if a signal delivery record precedes the exit;
- event order can support a mechanism claim, but cannot establish that the mechanism operated;
- neither can substitute for the other;
- the same observation must not be counted twice as independent evidence.

A circular pattern to reject:

1. The provider state is `stopping`, so an action is inferred.
2. The exit is observed while `stopping`, so ordering is inferred.
3. The same `stopping` state is treated as evidence that the action caused the exit.

That uses one mutable state as action identity, ordering, and mechanism. It is not three pieces of evidence.

The standard should require that action identity, event ordering, and mechanism operation each have independent support, or that the record explicitly says which are inferred and which are established.

## 5. Revised threshold model (analytical proposal)

This is a proposed refinement of the clarification memo's model, not a disposition.

### 5.1 Association quality

| Level | Evidence | Meaning |
|---|---|---|
| **A-E (event-order association)** | Action event and exit event are bound to the same instance, and event order is established by an authoritative per-instance sequence or equivalent. | The action event preceded the exit event. Still only association. |
| **A-O (observation-order association)** | Action and exit are observed in a recorded order, but event order is not established. | Observation order only; not sufficient for attribution. |
| **A-U (unresolved association)** | Same-instance binding or ordering is missing, concurrent, or contradictory. | No association; unresolved. |

### 5.2 Mechanism quality

| Level | Evidence | Meaning |
|---|---|---|
| **M-E (mechanism established)** | Specific pathway, evidence it operated, exit consistent, defeaters addressed. | Supports C-established with the other requirements. |
| **M-P (mechanism probable)** | Pathway and some operation evidence, but a material gap or unresolved defeater remains. | Supports C-probable at most; state the gap. |
| **M-U (mechanism unresolved or absent)** | Only capability, narrative, or correlation; or a defeater remains unexplained. | No causal attribution. |

### 5.3 Decision rule

- **C-established:** A-E + M-E + action outcome + exit fact + competing-explanation disposition + uncertainty statement.
- **C-probable:** A-E or A-O with M-P, or A-E with M-E but an unresolved non-material defeater; the record must state the limitation.
- **C-unresolved:** A-U, M-U, a material unexplained defeater, or contradictory evidence.

A-O alone never supports C-established, even if the mechanism narrative is plausible. M-P never becomes C-established by repetition or by a successful action result.

## 6. Effect on U5

U5 remains **OPEN**. The challenge does not close it and does not decide its public projection.

The current source still lacks:

- a durable action record with kind and target;
- a supervisor generation bound to action and exit;
- event-order evidence;
- a structured exit status;
- mechanism-operation evidence;
- competing-explanation handling.

The `stopping` state is still a mutable observation-time check. It cannot establish event order or mechanism operation. The four-step scenario remains the central hazard:

1. independent exit begins;
2. `Terminate` is issued;
3. nonzero exit is reported;
4. exit is observed while `stopping`.

All four can be true without action causation. U5 therefore cannot receive a unique, source-justified causal classification. It remains unresolved.

## 7. What this review does not change

This review does **not**:

- accept, reject, or defer OD-8;
- change OD-1 (deferred), OD-2/OD-3 (accepted), OD-4–OD-9 (unresolved), or U1–U8 (open);
- decide that U5 or any other case is action-caused;
- decide public projection, `crashed` semantics, quarantine, or compatibility taxonomy;
- modify the clarification memo, the OD-8 assessment, the decision register, P17, the corrected r3 response, the lifecycle decision, or contract rev3;
- create rev4, P18, a contract amendment, or an implementation;
- run tests, builds, benchmarks, or runtime experiments;
- claim that any proposed ordering or mechanism source currently exists;
- substitute for the owner's independent verification.

## 8. Verification note

The clarification memo is unchanged at:

`6dcc3b352b26b9d0ae4f6225ba97e4c51ed9a0446a9807cde0e350ec190ae552`

The OD-8 assessment is unchanged at:

`64447be9d166c3d7d9b6c09ce332180d5271168513f9f24b6aa91284472cb17e`

The decision register is unchanged at:

`4a2a24b908cbb382f05711d10750bc156e0abf7358eec1b7e98ef50116663a16` (244 lines)

The frozen artifact hashes remain:

- P17: `01d887c1…`
- Corrected r3 response: `60ca1cb5…`
- Lifecycle decision: `aa136bc9…`
- Contract rev3: `ed0da1a7…`

The helper aggregate digest remains `d0df8b28…`. This review made no source, test, configuration, or existing-document changes.

The owner's independent inspection and hash recomputation remain outstanding. This review does not claim to supply that independent verification.

**OD-8 remains UNRESOLVED. No register update or implementation is warranted.**
