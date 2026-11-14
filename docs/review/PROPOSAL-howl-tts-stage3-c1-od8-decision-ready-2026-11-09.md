# C1 / P17 — OD-8 Decision-Ready Proposal (Read-Only)

**Date:** 2026-11-09

**Status:** Proposal only. OD-8 remains **UNRESOLVED**. This document does not record a disposition, edit the decision register, or authorize implementation.

**Purpose:** Consolidate the association/attribution thresholds and the ordering/mechanism refinements into a single standard for owner decision, including an explicit outcome for genuinely unknowable event order.

**Inputs:**

- `docs/review/ASSESSMENT-howl-tts-stage3-c1-od8-causality-2026-11-09.md` — SHA-256 `64447be9…`.
- `docs/review/CLARIFICATION-howl-tts-stage3-c1-od8-attribution-thresholds-2026-11-09.md` — SHA-256 `6dcc3b35…`.
- `docs/review/REVIEW-howl-tts-stage3-c1-od8-ordering-mechanism-challenge-2026-11-09.md` — SHA-256 `0cc2ff7a…`.
- `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md` — SHA-256 `4a2a24b9…`, 244 lines.

## 1. Proposal in brief

OD-8 should adopt a **two-test standard with an explicit unresolved outcome**:

1. **Association test** — establishes that a specific action and an observed exit belong to the same process instance and are ordered. It never establishes causation.
2. **Attribution test** — adds action outcome, exit facts, evidence that a mechanism actually operated, and a disposition of competing explanations.
3. **Unresolved outcome** — applies when ordering is unknowable, a required element is missing or contradictory, or a material defeater remains. It is a legitimate result, not a failure to choose.

The standard distinguishes **event order** from **observation order**, and **mechanism capability** from **mechanism operation**. It also forbids using the same observation as simultaneous proof of action identity, ordering, and causation.

## 2. Why this proposal exists

The prior assessment established that action record + exit fact + ordering is association, not causation. The clarification separated the two thresholds. The challenge review found two remaining weaknesses:

- an action record appearing before an exit record does not prove the action's effect preceded the actual exit;
- a mechanism capable of producing an exit does not prove that it produced this exit.

This proposal consolidates the corrections and resolves the remaining question:

> What should the standard do when event order is genuinely unknowable?

The answer is: **do not force a causal conclusion.** Record an unresolved attribution and, if lifecycle closure is required, record the exit fact without converting architecture or necessity into evidence of causation.

## 3. Definitions

### 3.1 Event, observation, and record

| Term | Meaning |
|---|---|
| **Action event** | The action's operative occurrence for the purpose claimed: request, attempt, acceptance, effect delivery, or completion. For attribution, the relevant event is the operative effect. |
| **Exit event** | The process actually ending. |
| **Observation** | A component noticing an action or exit. |
| **Record** | A durable write describing an event or observation. |
| **Event order** | The order in which the action effect and exit actually occurred. |
| **Observation order** | The order in which a component observed or recorded them. |

Observation order is evidence about event order only when the observation mechanism is reliable for that purpose. It is not event order by definition.

### 3.2 Association quality

| Level | Meaning |
|---|---|
| **A-E (event-order association)** | The action effect and exit are bound to the same instance, and event order is established by an authoritative per-instance sequence or equivalent. |
| **A-O (observation-order association)** | Action and exit observations are ordered, but event order is not established. |
| **A-U (unresolved association)** | Same-instance binding or ordering is missing, concurrent, or contradictory. |

### 3.3 Ordering-source quality

A sequence number is only as good as its allocation point.

- A sequence allocated **at the event occurrence or acceptance point** can support event order.
- A sequence allocated **when a delayed observation is processed or persisted** merely reproduces observation order.
- A sequence shared by multiple writers can establish a total order only if all relevant events pass through that ordering authority. Otherwise the result is a partial order.
- Wall-clock timestamps are corroboration, not proof of event order. Clock skew, adjustment, resolution, asynchronous callbacks, and delayed persistence can all reverse apparent order.

### 3.4 Mechanism quality

| Level | Meaning |
|---|---|
| **M-E (mechanism established)** | A specific pathway is identified, evidence shows it operated on this instance, the exit is consistent with it, and material defeaters are absent. |
| **M-P (mechanism probable)** | The pathway and some operation evidence are present, but a material gap remains; confidence and the missing evidence are stated. |
| **M-U (mechanism unresolved or absent)** | Only capability, narrative, or correlation is present; or a material defeater remains unexplained. |

### 3.5 Attribution status

| Status | Meaning |
|---|---|
| **C-established** | Event order, mechanism operation, action outcome, exit fact, and competing-explanation disposition all meet the standard; no material defeater remains. |
| **C-probable** | The evidence supports attribution but a material uncertainty remains. It is not an established causal fact and must state its confidence and gap. |
| **C-unresolved** | Ordering is A-U, mechanism is M-U, a material defeater is unresolved, or the evidence is contradictory or insufficient. |

`C-unresolved` is an evidence status, not a provider state and not a public projection.

## 4. The proposed standard

### 4.1 Association test

Association requires all of the following:

| Element | Requirement |
|---|---|
| **A1. Action identity** | A specific action request or attempt is recorded, with its kind (terminate, kill, shutdown, cleanup, cancellation) and target instance. A generic `stopping` state is not an action identity. |
| **A2. Exit fact** | An actual process exit is observed, with the available exit status/error and observation context. A request, successful API return, or state change is not an exit fact. |
| **A3. Instance identity** | The action and exit are bound to the same authoritative process instance. A supervisor generation/incarnation token is the target concept under OD-2; PID is corroboration, not a replacement for generation. |
| **A4. Ordering** | The action event and exit event are ordered as A-E, A-O, or A-U. A-E requires event-order evidence. A-O is observation-only. A-U means no reliable order. |

**Rule:** A1–A4 establish association only. Association never establishes causation, and A-O or A-U never becomes C-established through repetition or a successful action result.

### 4.2 Attribution test

Attribution requires association plus:

| Element | Requirement |
|---|---|
| **C1. Action outcome** | The record distinguishes requested, accepted, attempted, completed, rejected, and timed out. A successful API return is not automatically acceptance or completion. |
| **C2. Exit fact quality** | The exit observation preserves the available exit code, signal, error, and observation context. A generic error is not a substitute for an exit status. |
| **C3. Mechanism operation** | Evidence that a specific mechanism capable of producing this exit actually operated on this instance, and that the observed exit is consistent with it. Capability alone is insufficient. |
| **C4. Material defeater** | At least one material defeater is identified, and its status is recorded as present, absent, or unresolved. A defeater is material if its presence would defeat the attribution. |
| **C5. Competing explanations** | A pre-existing exit, independent failure, fatal observation, shutdown-time failure, or second action is recorded and either excluded by evidence or left explicitly unresolved. |
| **C6. Uncertainty** | The record states whether the result is C-established, C-probable, or C-unresolved, and what evidence supports that status. |

**Rule:** C-established requires A-E, mechanism operation (M-E), no material defeater, and a disposition of competing explanations. C-probable is permitted only with an explicit confidence statement and stated gap. C-unresolved applies otherwise.

### 4.3 No circular evidence

The same observation must not be counted as independent support for more than one required element.

In particular, `providerState == stopping` cannot simultaneously serve as:

1. proof that a specific action occurred;
2. proof that the action preceded the exit;
3. proof that the action caused the exit.

It is at most one state observation. Action identity, ordering, and mechanism operation need independent support, or must be explicitly labeled inferred. A terminal decision requirement is not evidence and cannot fill any of these roles.

## 5. Genuinely unknowable event order

### 5.1 Definition

Event order is **genuinely unknowable** when:

- no authoritative event-order record exists;
- the available observations are concurrent, overlapping, or recorded at points that do not correspond to event occurrence; and
- no independent mechanism evidence establishes which event produced the exit.

This is distinct from **currently unknown**, where the evidence is missing but could still be obtained. Both are unresolved now; only currently-unknown cases can later become A-E through new evidence.

### 5.2 Required outcome

| Situation | Association | Attribution | Required record |
|---|---|---|---|
| Event order established by an authoritative source. | A-E | C-established, C-probable, or C-unresolved depending on mechanism and defeaters. | Record the ordering source and its allocation point. |
| Only observation order established. | A-O | C-established is impossible. C-probable only if independent mechanism evidence supports operation; otherwise C-unresolved. | State that event order is unknown and why. |
| No reliable order; concurrent or unknown. | A-U | C-unresolved. C-probable and C-established are unavailable without A-O or A-E. | Record the unresolved ordering and the missing evidence. |
| Genuinely unknowable and no independent mechanism evidence. | A-U | C-unresolved. | Do not infer a cause from terminal necessity, class order, or elapsed time. |
| Terminal lifecycle decision is required before attribution can be resolved. | A-U or A-O as applicable. | C-unresolved may coexist with terminal closure. | Record the exit fact; record attribution as unresolved; do not record action-caused by default. |

### 5.3 Rules for unknowable order

1. **A-U never supports C-established.** No amount of narrative or architectural urgency changes this.
2. **A-O never supports C-established.** Observation order is not event order.
3. **Genuinely unknowable order is C-unresolved.** A-U never supports C-probable or C-established. C-probable requires at least A-O; C-established requires A-E.
4. **Terminal necessity is not evidence.** If the architecture requires a terminal lifecycle decision, that requirement may justify recording the exit fact and closing the lifecycle record; it does not justify asserting causation.
5. **Late evidence may resolve a currently-unknown order.** If an authoritative event-order record arrives, A-U may become A-E. If the order is genuinely unknowable, it remains unresolved.
6. **Do not default to `crashed`, `absent`, or `terminated_by_action`.** Those are projections or internal dispositions; their selection is outside this standard and remains subject to OD-4b/OD-5.

### 5.4 Example: asynchronous exit and action return

1. The OS reports the exit asynchronously.
2. The action API returns around the same time.
3. Neither timestamp nor record-time sequence establishes whether the action effect preceded the exit.
4. No independent signal-delivery or process-audit evidence is available.

**Result:** A-U; C-unresolved. A terminal lifecycle record may say the instance ended and record the exit status. It must not say the action caused the exit. If later evidence establishes the action effect preceded the exit, the attribution may be revisited; the immutable terminal lifecycle decision is not rewritten.

## 6. Terminal decisions and late evidence

### 6.1 Two separate questions

The standard distinguishes:

- **Lifecycle closure:** did this process instance end, and what exit facts are available?
- **Causal attribution:** did this action cause that exit?

A terminal lifecycle decision may be made on the first question without answering the second. The second may remain C-unresolved after closure.

This proposal does not decide the representation of lifecycle closure or the public projection. It states the evidentiary rule: **a requirement to close the lifecycle is not evidence of causation.**

### 6.2 OD-3 interaction

OD-3's accepted principle is an immutable first terminal decision per process instance, with later observations retained. This proposal preserves that principle:

- the terminal decision is not rewritten;
- late evidence remains recordable;
- late evidence may strengthen, weaken, or leave the attribution unresolved;
- whether late evidence may change a published projection remains an OD-3/OD-4b question, not decided here.

If a terminal decision and an attribution judgment are represented together, the attribution must not be treated as established merely because the lifecycle decision is terminal. If they are represented separately, the separation must be explicit. The representation choice is outside this proposal.

### 6.3 OD-1 interaction

OD-1 remains deferred. This proposal defines evidence requirements, not evidence ownership. An authoritative recorder, publisher, and synchronization mechanism are still required before the standard is operational; their absence is not a reason to weaken the standard or to infer causation.

## 7. Decision table

| Ordering | Mechanism | Defeater | Result |
|---|---|---|---|
| A-E | M-E | No material defeater; competing explanations addressed. | **C-established.** |
| A-E | M-P | Material gap stated; no material defeater. | **C-probable**, with confidence and gap. |
| A-E | M-E | Material defeater unresolved. | **C-unresolved.** |
| A-O | M-E | No material defeater; event order unknown. | **C-probable at most**; never C-established. |
| A-O | M-P | Event order unknown; mechanism probable. | **C-unresolved by default**; C-probable only with independent operation evidence and explicit order uncertainty. |
| A-U | Any | Any. | **C-unresolved**. Independent evidence may reclassify the ordering as A-O or A-E; until then, no attribution label. |
| Any | M-U | Any. | **C-unresolved.** |
| Any | Any | Material defeater present. | **No action-caused attribution**; record the competing explanation. |

A `C-probable` result is not a public provider state and must not be presented as an established causal fact.

## 8. Owner decision points

This proposal is ready for an owner decision on the following points. It does not decide them.

### 8.1 Adopt the two-test structure?

**Recommended:** yes. Association and attribution answer different questions and should not be collapsed.

### 8.2 Allow C-probable, or require all non-established cases to be C-unresolved?

**Option 1 (recommended):** allow C-probable with an explicit confidence statement and stated gap; never present it as C-established.
**Option 2 (stricter):** only C-established or C-unresolved; no probabilistic attribution label.

Either option is compatible with the no-forced-causation rule. The owner should choose whether a probabilistic attribution is useful enough to retain as a distinct evidence status.

### 8.3 Allow A-O to support C-probable?

**Option 1 (recommended):** A-O may support C-probable only when independent mechanism evidence supports operation and the record explicitly includes order uncertainty.
**Option 2 (stricter):** A-O always yields C-unresolved, even with mechanism evidence.

The stricter option avoids any causal label when event order is unknown. The recommended option preserves useful probabilistic evidence without promoting it to C-established.

### 8.4 Require a material defeater?

**Recommended:** yes. Without a defeater, a mechanism claim is not testable. The record must identify at least one material defeater and state whether it is present, absent, or unresolved.

### 8.5 Confirm terminal closure does not imply causation?

**Recommended:** yes. A terminal lifecycle decision may record the exit fact; it must not fill an unresolved causal attribution by default. The representation and public projection remain separate decisions.

### 8.6 Keep U5 open?

**Recommended:** yes. This proposal does not classify U5. U5 remains an application case to be evaluated against the eventual standard.

## 9. Effect on U5

U5 remains **OPEN**. The current source still lacks:

- durable action identity and kind (`Kill` is an alias of `Terminate`);
- supervisor generation bound to action and exit;
- event-order evidence (no sequence; `stopping` is an observation-time state);
- structured exit status (only `waitErr`);
- mechanism-operation evidence;
- material-defeater and competing-explanation records.

The four-step scenario still applies:

1. an independent exit begins;
2. `Terminate` is issued;
3. a nonzero exit is reported;
4. the exit is observed while `stopping`.

All four can be true without action causation. U5 therefore cannot receive a unique, source-justified causal classification. It remains unresolved, and no public projection is decided.

## 10. What this proposal does not decide

This proposal does **not**:

- record an OD-8 disposition or declare the standard accepted;
- change OD-1 (deferred), OD-2/OD-3 (accepted), OD-4–OD-9 (unresolved), or U1–U8 (open);
- decide that U5 or any other case is action-caused;
- decide `crashed` semantics, public projection, quarantine, or compatibility taxonomy;
- decide evidence-record ownership, publisher authorization, synchronization, or late-evidence projection effects;
- modify the clarification memo, challenge review, assessment, decision register, P17, corrected r3 response, lifecycle decision, or contract rev3;
- create rev4, P18, a contract amendment, or an implementation;
- run tests, builds, benchmarks, or runtime experiments;
- claim that any proposed ordering source or mechanism record currently exists;
- substitute for the owner's independent verification.

## 11. Verification note

The prior artifacts are unchanged at:

- Clarification memo: `6dcc3b352b26b9d0ae4f6225ba97e4c51ed9a0446a9807cde0e350ec190ae552`
- Challenge review: `0cc2ff7a636b747cf7e9567a7b714fac7ae8ecea2b0fab7900b084f1d278506e`
- OD-8 assessment: `64447be9d166c3d7d9b6c09ce332180d5271168513f9f24b6aa91284472cb17e`
- Decision register: `4a2a24b908cbb382f05711d10750bc156e0abf7358eec1b7e98ef50116663a16` (244 lines)

Frozen artifacts remain:

- P17: `01d887c1…`
- Corrected r3 response: `60ca1cb5…`
- Lifecycle decision: `aa136bc9…`
- Contract rev3: `ed0da1a7…`

The helper aggregate digest remains `d0df8b28…`. No source, test, configuration, or existing document was modified.

The owner's independent inspection and hash recomputation remain outstanding. This proposal does not claim to supply that independent verification.

**OD-8 remains UNRESOLVED. No register update, contract change, or implementation is warranted.**
