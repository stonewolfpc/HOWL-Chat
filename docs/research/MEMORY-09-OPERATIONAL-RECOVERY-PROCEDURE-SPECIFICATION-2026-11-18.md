# MEMORY-09 — Operational Recovery: Procedure Specification

**Status: SPECIFICATION — DESCRIPTIVE ONLY. NOT IMPLEMENTED. NOT TESTED. NOT AUTHORIZED FOR EXECUTION.**

- Authored under §4 of `MEMORY-09-OPERATIONAL-RECOVERY-DECISION-STATEMENT-2026-11-18.md`, which
  authorizes **specification authoring and review only**.
- Introduces **no new mechanism**. It describes the boundary of the mechanisms that already exist.
- Does **not** resolve E (missing-record repair) or F (new-root initialization). Both remain open
  and appear below as **stop points**, not as solved branches.
- Does **not** change Finding 2, any finding, any total, or any trust-model decision.
- **Nothing here has been run, rehearsed, or verified** against a real or synthetic chain.

## Provenance

| Field | Value |
|---|---|
| Document class | Specification — descriptive, unexecuted |
| Derived from | The reviewed operational-recovery analysis (Pass 1 and Pass 2 complete) |
| Transcribed vs added | §3–§8 restate the analysis. §9–§13 are **author proposals** structuring that material into a procedure; they add no requirements |
| Implementation authority | **None** |
| Verification state | No execution, no rehearsal, no harness, no corpus |
| Canonical record impact | **None** |

## 1. The situation this addresses

Acceptance is **blocked** because chain authority is contested: the validator cannot establish
which presented continuation of a known head is authoritative, or the presentation it received
cannot be shown to be complete.

This is **not** a description of normal operation. In normal operation the chain provides internal
consistency, citation-topology ordering, and refusal of visibly ambiguous presentations. Nothing
below applies until a pending acceptance is actually blocked.

Recovery here means **resuming acceptance after a stall**, not restoring data, repairing damage, or
proving that history is authentic.

## 2. Three layers

| Layer | Does | Does not |
|---|---|---|
| **Mechanical** | Observes, filters, classifies, refuses, preserves | Adjudicate, select, or authorize |
| **Authority** | Declares successor, abandonment, or resumption, with scope and provenance | Prove the declaration's own authorship |
| **Evidence** | Preserves observations and declarations as what they are | Convert either into a stronger claim than its provenance supports |

A declaration is recorded **as** a declaration. There is no layer at which the instrument becomes
the authority.

## 3. Step 1 — Detection

Only an observation that blocks a **pending acceptance** is actionable. A quiescent unresolved chain
blocks nothing and therefore triggers nothing.

| Observation | Character | Standing |
|---|---|---|
| Referenced artifact absent; reference unresolvable | gap | `incomplete` — proposed classification |
| Referenced artifact present, digest ≠ recorded | conflict | contradiction family — subtype coverage not yet justified |
| Competing successors, co-presented | conflict | existing |
| Multiple heads, no valid path | conflict | existing |
| Cycle at the head | structural conflict | existing |
| Presentation asymmetry — two parties report different heads | **not detectable by one validator** | social trigger |

Presentation asymmetry is **not a malfunction**. It is a limit on the evidence presented. A single
validator cannot detect it, so it cannot self-trigger.

*Note:* whether quiescence-triggerlessness merely restates the liveness limitation or introduces a
new rule is **unresolved** (reviewer question, not a decision). This specification does not settle it.

## 4. Step 2 — Preservation

The collection is written **outside** the contested presentation. Adding a manifest inside the
contested directory would mutate the state being preserved and could alter the apparent head set.

Recorded:

- the presented directory or snapshot identity, with method and time;
- every artifact with its digest;
- the observed head set;
- which validator observed which presentation;
- refusals and trigger observations, with time and source;
- the terminal record, if preserved, **as a claim only — never as an adjudicator**.

Record **all observed branches**, not "both". Asymmetry means the known count cannot be assumed
complete.

Preservation and repair are strictly separate. Locating a previously omitted artifact may alter the
committed evidence set under the same `run_id`, and no testable discriminator exists between that and
an impermissible evidence-set change. See **stop point S2**.

## 5. Step 3 — Continuity test

The anchor is an **exclusion filter**. It excludes non-continuations; it does not select an
authorized continuation.

| Condition | Mechanical result |
|---|---|
| Candidate makes **no claim** of continuity with H | **Excluded**. No continuity defect is asserted or implied for that candidate |
| Candidate **claims** continuity but cannot demonstrate it | **`incomplete`** — refusal and escalation |
| H **absent from all** presented candidates | **No global conclusion** — continued refusal or escalation where continuity is required |

The third row is not a fourth kind of evidence defect. It is the limit of what the presentation
establishes: the instrument cannot certify that the candidate set is complete, so an empty filter
result is evidence **about the presentation**, never proof that no continuation exists.

Where exactly one candidate survives, it is the sole *presented* continuation of H **within this
presentation** — uncontested, and **not authorized**. Owner acceptance is still required.

## 6. Step 4 — Freshness bound

With attested heads H₁ < … < Hₙ, the test partitions on Hₙ and leaves the interval (Hₙ, head]
unanchored. Divergence inside that interval is invisible.

A fully attested sequence pins history only **through its last attested point**. Closing the interval
requires either an attested current head or an explicitly closed sequence — neither of which exists
under the current trust-model decisions. Freshness belongs to (c), which remains **deferred**.

This is a statement of the bound, not a proposal to close it.

## 7. Step 5 — Resumption record

Three separately recorded facts, with distinct provenance. **No synthesised `recovery_succeeded`
verdict is produced.**

| Field | Provenance | Establishes |
|---|---|---|
| Traversal result, scoped to the presentation examined | mechanical | no detectable conflict within this presentation |
| Owner designation | authority declaration | the history accepted for continued use |
| Non-existence of other branches | **unestablished** | not establishable by this instrument |

The designation names the **specific head by digest** — never a directory path or branch label —
with scope and time. Labels are mutable; only digests are identities.

## 8. Step 6 — Abandonment, where the owner takes it

Abandonment has **no mechanical trigger** and cannot be inferred by the instrument. It belongs
exclusively to the authority layer.

- The owner may **declare** the old chain abandoned.
- The declaration records **who, when, scope, reason**.
- It does **not** claim the instrument proved no continuation existed.
- The new root starts a **new chain** and inherits **no** old acceptance authority.
- The declaration's authorship is **asserted, not independently authenticated**.

Citation boundary, three parts: prior records remain citeable as surviving documents or historical
claims; they establish no authenticated ordering or continuity through the abandoned chain; the new
root inherits none of their authority. **Content survives; ordering and continuity do not.**

Abandonment is **not** D8-U4-iii succession. A successor under U4-iii continues under the original
run's authority; abandonment ends that chain and starts another carrying none of it.

Abandonment does not by itself answer what anchors the new root. See **stop point S5**.

## 9. Stopping conditions

The procedure stops and escalates at each of the following. These are the conditions under which the
instrument must **refuse to proceed**, stated so that refusal is the specified behaviour rather than
an improvised one.

| # | Condition | Why it stops |
|---|---|---|
| **S1** | Presentation asymmetry is reported by two parties | Not detectable by one validator; the trigger is necessarily social |
| **S2** | A previously omitted artifact is located and argued to be recoverable | E — no testable discriminator exists between mechanical repair and an impermissible evidence-set change |
| **S3** | H is absent from every presented candidate | The presentation cannot establish completeness; no global conclusion follows |
| **S4** | Competing successors, multiple heads, or a head cycle | Visible ambiguity is refused, preserved, and escalated. Digest-mismatch and gap-`incomplete` refusals are specified in §3/§5 and are not restated here |
| **S5** | A decision relies on abandonment, or on what anchors a new root | F — new-root initialization is unresolved; a new root is not self-authorizing |
| **S6** | Any request to assert operational safety or a tested procedure | Neither is established; the tested-procedure limit applies |
| **S7** | Any query requiring the completeness of the branch set | Branch non-existence is unestablishable by this instrument |

Stopping is the specified successful stop — **not** a successful recovery. It preserves the
distinction between "no conflict was detected" and "no conflict exists."

## 10. What the instrument cannot know

| Question | Answer |
|---|---|
| Is this presentation complete? | No — never established |
| Are there other branches? | Unknown and not establishable |
| Is this artifact earlier than that one? | No — hashes establish neither chronology nor authorship |
| Is the root authentic? | No — the root is replaceable |
| Was the chain rewritten with recomputed digests? | Observationally indistinguishable from within |
| Is the acceptance safe? | Undemonstrated |

## 11. Owner-only decisions

Not delegable to the instrument, and not inferable from any observation above:

1. Which head is accepted for continued use — named by digest.
2. Whether an old chain is abandoned.
3. Whether to reconsider (b) signatures or (c) external head witnessing.
4. Whether a previously omitted artifact is a mechanical repair (E).
5. What anchors a new root (F).
6. Whether any recovery procedure is safe enough to authorize for use.

## 12. Boundaries this specification deliberately leaves open

| Boundary | State here |
|---|---|
| **E** — D8-U4-iii / §7 item 11 discriminator | Unresolved. Appears as stop point S2. Not answered |
| **F** — new-root initialization | Unresolved. Appears as stop point S5. Not answered |
| **(b) signatures / (c) witnessing** | Remain deferred. Freshness is stated as a bound, not closed |
| **D8-U2 subtype coverage** | Proposed, unjustified; the digest-mismatch row carries that status |
| **Quiescence provenance** | Reviewer question; restatement-versus-new-rule not settled |
| **Artifact governance** across specification, authorization, and evidence | Out of scope here |

Where this specification needed an answer and none exists, it stops rather than supplies one.

## 13. Proportionality

This is a local application's internal memory-integrity design. The specification deliberately adds
**none** of the following:

- key management, signature infrastructure, or a certificate hierarchy;
- external witnesses, timestamping services, or third-party dependencies;
- an in-system adjudicator, arbiter, or voting mechanism;
- new storage, new services, or new trust boundaries;
- a governance body, review board, or approval workflow.

The instruments it describes are a digest-linked chain, a refusal path, an owner declaration, and a
preserved record. Recovery is an owner act recorded honestly, not a subsystem.

## 14. Verification state

| Check | State |
|---|---|
| Implemented | **No** |
| Executed or rehearsed | **No** |
| Corpus prepared | **No** |
| Tested against a real contested history | **No** |
| Safety demonstrated | **No** |

This document describes an intended procedure. It is not evidence that the procedure works, and no
statement here may be cited as such.

---

## Status

Descriptive specification only. Finding 2 remains `false_closure`, partially resolved, open on
operational recovery. E and F remain open. (b) and (c) remain deferred. No finding, total,
classification, or trust-model decision is altered. Nothing is implemented, tested, or authorized
for execution.
