# MEMORY-09 — OD-b Option Evaluation Against the Revised Criteria

## L1, L2 and L3 assessed; the choice that remains stated, not made

**Status:** analysis record. **No selection in this record.** OD-b was subsequently **settled by the owner:
L2** (2026-11-20), recorded in `MEMORY-09-OWNER-DECISION-OD-B-L2-2026-11-20.md`; §11 states what that changes
here and what it does not. The record itself still evaluates and stops. Read-only: no schema, no field list, no
canonical edit, no test, no ledger entry, no probe.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Criteria | ODB-64 / ODB-65 / ODB-66 (brief revision 3, §5.2) |
| Question answered | *What does each option cost, what does it buy, and how does it fail silently?* |
| Question **not** answered here | *Which one to select* — answered by the owner **2026-11-20: L2** (§11) |
| Decision | **OD-b = L2**, owner, 2026-11-20 — recorded elsewhere; §11 (OBE-40…OBE-47) |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Method, and one correction to the framing

**OBE-01 (Observed — criteria as given).** Three questions, one per option: L1's naming/vocabulary cost
(ODB-64), L2's fold/reason/closure completeness (ODB-65), L3's determinism versus unintended authority
(ODB-66). They are **not commensurable** — each measures a different kind of cost — so this record does not add
them up. It reports, per option, an **obligation** and a **silent-failure mode**.

**OBE-02 (Inferred — why "silent-failure mode" is the right common denominator).** The options are not being
compared on capability: all three can express the five documented facts (ODB-40, OBE-20). They differ in **what
goes wrong when something is omitted**, which is the property an owner actually lives with, because omission is
the normal state of a partially-populated stream.

**OBE-03 (Inferred — deferral is not avoidance).** A correction to how L1's cost is usually stated. L2 and L3
**defer** the `active`/`established` naming to a fold or a derivation; they do not escape it. The projection
still emits one name, so the vocabulary ruling (K-R12) is needed **either way**. L1's cost is therefore
**timing** (it must be settled before representation), not magnitude. On this criterion the three are closer than
ODB-64 implies.

**OBE-04 (Inferred — ordering is not an L1 penalty).** An ordered state sequence needs an ordering key exactly as
a relation set and an act set do. That requirement is **common to all three** and is the subject of the separate
ordering-key comparison, not a discriminator here. Recording this prevents L1 from appearing to lose on a cost
every option pays.

**OBE-05 (Observed — hybrids not assessed).** A combination of stances would be a **fourth option**, not a
refinement of these three. None is evaluated here, and none is assumed available (brief §12).

---

## 2. Criterion decomposition

**OBE-06 (Inferred — sub-criteria, so that "pass" means something specific).**

| ID | Sub-criterion | Applies to |
|---|---|---|
| S1 | Can the five documented facts (ODB-40) be expressed? | all |
| S2 | Is there a single source of truth for status? | all |
| S3 | Is the terminal `provisional` (rumour) representable without inference? | all |
| S4 | Does retraction retain prior acquisitions (`:616`)? | all |
| S5 | Is the projection rule of K-R12 satisfied structurally? | all |
| S6 | Does an omitted fact fail **visibly** or **silently**? | all |
| S7 | Does the option require a decision from an open owner question as a precondition? | all |
| S8 | Can the `active`/`established` naming be deferred? | all |
| S9 | Does the option resist absorbing forgetting (DN-4)? | all |

**OBE-07 (Inferred).** S1–S4 are the documented facts; S5–S9 are the properties that decide what the owner
inherits. An option that fails S1–S4 is not viable; S5–S9 is where the real trade sits.

---

## 3. L1 — State-naming

**OBE-08 (Calculated — assessment).**

| Sub-criterion | Result | Basis |
|---|---|---|
| S1 five facts | **pass** | a state name per condition; `corrected` differs from a re-derived fact (ODB-32) |
| S2 single source | **partial** | state and its ground can disagree; two sources exist for one condition |
| S3 terminal provisional | **pass, best** | the state is stored, so no closure argument is needed (ECI-25 avoided) |
| S4 retraction retention | **pass** | the state names the target only; prior facts untouched |
| S5 K-R12 structure | **partial** | a stored state is not a projection — it must be derived from, or asserted in, the stream |
| S6 silent failure | **assertion without warrant** (OBE-11) | — |
| S7 open-decision precondition | **weak** — the vocabulary must be fixed first | K-R12 naming is open |
| S8 deferral of naming | **no** — this is L1's characteristic cost | OBE-03 |
| S9 forgetting resistance | **exposed, not forced** | a status vocabulary invites absorption (CLV-38) |

**OBE-09 (Inferred — what L1 buys).** No fold, therefore **no default**, therefore **no phantom provisional**
(ECI-27) and **no closure mechanism** (ECI-25). The terminal rumour is simply a record whose last state is
`provisional`. On the two failure modes that ECI identified as L2/L3 obligations, L1 pays nothing.

**OBE-10 (Inferred — L1's obligation).** One thing: **keep forgetting out of the state vocabulary** while
explaining why. This is a documentation burden with a live failure: the easiest way to record "forgotten" is a
status, and DN-4 forbids it.

**OBE-11 (Inferred — L1's silent-failure mode).** **A state the record's evidence cannot support.** Storing a
condition makes it possible to state a condition without recording what produced it. Nothing in the structure then
distinguishes "was promoted on witnessing" from "was marked active." The state is not wrong; it is **unwarranted**,
and unwarranted states look exactly like warranted ones.

**OBE-12 (Inferred — the honest size of this).** It is not a correctness bug; every representation can be
misused. It is a **verification** cost: L1 makes status checkable only against the lifecycle stream's own
provenance, which under OD-a Option B is a separate structure that may hold only states. The asymmetry is that L1's
failure surfaces at **audit** time, while L2's and L3's surface at **read** time.

---

## 4. L2 — Transition-only

**OBE-13 (Calculated — assessment).**

| Sub-criterion | Result | Basis |
|---|---|---|
| S1 five facts | **pass** | each transition is a distinguishable change; ordering supplies sequence |
| S2 single source | **pass, best** | nothing stored but the transitions; status is always derived |
| S3 terminal provisional | **partial** — needs closure | a rumour and a truncated history fold identically (OBE-15) |
| S4 retraction retention | **pass, best** | an append-only relation set cannot un-happen anything (`:616`) |
| S5 K-R12 structure | **pass, best** | status is a projection by construction —
 exactly what `:111` states |
| S6 silent failure | **truncation indistinguishable from completion** (OBE-15) | — |
| S7 open-decision precondition | **mild** | carrying a reason reaches toward the ground (OBE-14) |
| S8 deferral of naming | **yes** — deferred to the fold | OBE-03 |
| S9 forgetting resistance | **neutral** | no vocabulary to absorb it |

**OBE-14 (Inferred — L2's obligation, restated precisely).** Two, and the first has a hidden second step.

1. **Retain the reason for a transition.** A relation records *that*, not *why* (CLV-27). Retaining it means a
   ground reference inside the transition — which is where a cause can begin to authorise a status (CLV-28). L2
   may decline the reference, at the cost that "why" is unrecoverable later.
2. **Close the stream.** ECI-25. This is not a relation and not a state. **If closure must be expressed as an
   additional kind of fact, L2 has acquired a second structural element** — narrowing the distance to L1 without
   removing it, since a closure marker asserts *that the history ends here*, not *what the condition is*.

**OBE-15 (Inferred — L2's silent-failure mode).** **Structural truncation reads as completion.** A record whose
history was lost folds to `provisional`, which is precisely the correct answer for a rumour that never exits. The
fold cannot distinguish them, so the failure produces a **plausible, confident `provisional`** — the worst shape a
failure can take in a memory system, because it looks like a finding. Obligation 2 exists exactly to prevent this,
which is why L2's closure burden is substantive rather than administrative (ODB-65).

**OBE-16 (Inferred — the §8.1 interaction, scoped).** Under L2's default, a formation that legitimately produced
`active` without an entry fact folds to `provisional` — a wrong state from a correct default (EKR-22). This is
**not** an argument against L2: it is an argument that a fold-based selection **inherits** the §8.1 question as a
downstream obligation (EKR-24). It does not settle it.

---

## 5. L3 — Act-only

**OBE-17 (Calculated — assessment).**

| Sub-criterion | Result | Basis |
|---|---|---|
| S1 five facts | **pass** | an act per change; the act's kind distinguishes corrected from superseded |
| S2 single source | **pass** | derivation from acts |
| S3 terminal provisional | **partial** — no exit act exists, so absence must mean terminal | same ambiguity as L2 unless closure is added |
| S4 retraction retention | **pass** | an act cannot be undone by a later act; retention is structural |
| S5 K-R12 structure | **pass** | status derived, never stored |
| S6 silent failure | **recording confers** (OBE-19) | — |
| S7 open-decision precondition | **strong** — validity is entangled | OBE-18 |
| S8 deferral of naming | **yes** | — |
| S9 forgetting resistance | **neutral** | — |

**OBE-18 (Inferred — L3's obligation, and why it is the heaviest).** Deterministic derivation must be **total**
across every kind, including the kind that never exits, and it must not accept an act merely because it was
recorded. That second half is a **validity** question: which recordings count. ND-1 (who may advance) is open, so
L3's derivation cannot be specified without touching an open owner decision in a way L1 does not and L2 only
mildly does (OBE-14.1).

**OBE-19 (Inferred — L3's silent-failure mode).** **Recording an act confers the status.** If the derivation
reads "an act of kind K is present → state S," then the ability to append an act is the ability to change status.
The separation failure here runs in the opposite direction from CLV-26: not an omission, but an **inversion** —
the representational element acquires the authority it was supposed to describe.

**OBE-20 (Inferred — the honest counterweight).** L3 has no duplication and no default, so it cannot produce a
phantom provisional by defaulting; it produces one only if derivation treats absence of an act as a state, which
is a specified choice rather than an accident. That is a real advantage over L2, and it is the one place where L3
is strictly cleaner than the fold.

---

## 6. Where the three actually differ

**OBE-21 (Calculated — the comparison that discriminates).**

| | L1 | L2 | L3 |
|---|---|---|---|
| **Characteristic obligation** | keep forgetting out of the vocabulary | closure, plus a decision on carrying the reason | totality of derivation, plus validity |
| **Silent-failure mode** | unwarranted state | truncation reads as completion | recording confers |
| **Surfaces at** | audit | read | write |
| **Open decision it touches** | K-R12 naming (mild) | the ground, if the reason is carried (mild) | ND-1 validity (strong) |
| **Buys** | no fold, no closure, no default | single truth, structural retention, K-R12 fit | no default, no duplication, structural retention |

**OBE-22 (Inferred — S6 is the real axis).** The three failure modes differ in **when** and **to whom** they
appear. L1's is invisible until someone audits provenance; L2's produces a confident wrong state at read time;
L3's lets write access become status authority. These are not comparable on a single scale, and choosing among
them is a judgement about which failure this system can least afford — a product/risk judgement, not a
representational one.

**OBE-23 (Inferred — no dominance, confirmed).** No option passes every sub-criterion better than the others.
L2 wins S2, S4 and S5; L1 wins S3 and S8's opposite (certainty over deferral); L3 ties on S2 and S4 and wins
nothing uniquely except the default-free property (OBE-20). This is the strongest available statement that **the
choice is a trade, not a discovery**.

**OBE-24 (Observed — C′ and OD-g remain equal).** None of the three has a home for the verification dimension
(C′), and none preempts OD-g. Equal on both; neither is a discriminator (brief §10, §11).

**OBE-25 (Inferred — DN-4 exposure is asymmetric but small).** Only L1 has a vocabulary that could absorb
forgetting, and it can refuse. So DN-4 does not decide between the options; it constrains **how L1 would be
implemented** if selected.

---

## 7. What would falsify an option

**OBE-26 (Inferred — stated as tests of viability, not as preferences).**

| Option | Falsifier | Consequence if it holds |
|---|---|---|
| L1 | a lifecycle concept that **must** be a status value and cannot be a name outside the five | L1's vocabulary cannot cover the progression |
| L2 | closure cannot be specified without a successor relation — i.e. "no successor" cannot be distinguished from "unknown" | L2's S3 stays partial and its silent failure is unfixable |
| L3 | no **total** derivation exists over the kinds (belief is currently the candidate — EKR-13) | L3 cannot be deterministic |

**OBE-27 (Observed — a genuine asymmetry in falsifier status).** L2's and L3's falsifiers are **live open
questions** (EKR-13 leaves belief unresolved; closure is unspecified). L1's falsifier is **not live**: no canon
text requires any concept beyond the five names to be a status. So L1 is the only option whose viability is
currently unthreatened by an open question — a fact about the state of knowledge, **not** a preference, and not a
recommendation.

**OBE-28 (Inferred — how to read OBE-27).** It means a selection of L2 or L3 carries a **conditional**: it is
contingent on closing the falsifier. It does not mean L1 is better; L1's costs are real, they are simply not
gated on anything unknown.

---

## 8. The choice, stated without making it

**OBE-29 (Inferred — the question in one sentence).** *Which of these three failures can this system least
afford: an unwarranted state (audit-time), a truncated history that reads as complete (read-time), or a write that
confers status (write-time)?*

**OBE-30 (Inferred — the obligations that follow, per answer).**

- Choose **L1** → accept a documentation obligation (forgetting stays outside the vocabulary) and an audit
  obligation (states checked against provenance).
- Choose **L2** → accept a closure mechanism, probably a second kind of fact, and a decision on whether the reason
  for a transition is retained.
- Choose **L3** → accept derivation totality across all six kinds and a validity rule that will touch ND-1.

**OBE-37 (Observed — form refined by the owner).** OBE-29's wording is superseded in **form** (not in substance) by:
> *"Which failure mode is least tolerable, considering both its likelihood and the difficulty of detecting and
> recovering from it?"*

The refinement matters because exposure is not equal (CEV-40). OBE-29's substance — that the choice is a risk
judgement, not a representational one — stands unchanged.

**OBE-38 (Observed — caution carried).** OBE-27 must not be read as "L1 is safest overall." Certainty about
prerequisites and safety of the resulting choice are different claims; L1's vocabulary and state-representation
costs are untouched by it.

**OBE-39 (Observed — the two conditional viability questions now have evidence).** OBE-36 and OBE-18's ND-1
entanglement are examined in `...CLOSURE-AND-VALIDITY-EVIDENCE...`: closure is **required in effect and absent in
text** (CEV-18), and validity is **stated and unenforceable** (CEV-32). Neither is resolved; both are better posed.

**OBE-31 (Observed — no selection here).** This record does not select, rank, or weight. OBE-23 is the closest it
comes, and OBE-23 says there is no dominance to report. **Updated 2026-11-20:** the owner selected L2; that choice
is recorded in §11, not here, and this record still neither selects, ranks nor weights.

---

## 9. Separation from the other open questions

**OBE-32 (Observed).** §8.1 is **not** used as a discriminator here (OBE-16). Entry is constrained and identical
across the three options (ODB-63), so nothing in this evaluation depends on it.

**OBE-33 (Observed).** `active`/`established` remains a **vocabulary** question; this record only observes that
it is deferred by L2/L3 and settled earlier by L1 (OBE-03), not that either is preferable.

**OBE-34 (Observed).** C′, OD-g, ND-1, DN-4 and the domain identity question are untouched, and no field, event
type, or mechanism is proposed.

---

## 10. Limits

**OBE-35 (Observed).** Basis: brief revision 3 §5.2/§6, ECI-21/25/27/31, CLV-26/27/28/38, EKR-22/24, r33 K-R12
`:111`, §11 `:616`. Not inspected: any implementation's lifecycle handling; whether a closure marker can be
expressed as a relation in the existing stream vocabulary (this would materially soften OBE-14.2 and is
**Unestablished**); the cost of totality in practice.

**OBE-36 (Inferred — the weakest point in this record).** OBE-14.2's claim that a closure marker is "a second kind
of fact" assumes relations and states exhaust the vocabulary. If a stream can express "no further entries by
design," L2's obligation collapses to a marker and its distance from L1 widens again. This is the one place where
the comparison could shift materially on new text.

---

## 11. The owner's selection, recorded 2026-11-20

### 11.1 What was selected, and by whom

**OBE-40 (Observed — owner).** The owner selected **L2 — transition-only**: lifecycle history is represented
through transitions, and current lifecycle state is derived as a projection of that authoritative history. **L1 and
L3 were not selected.** The decision statement, rationale and failure-risk ordering are recorded in
`MEMORY-09-OWNER-DECISION-OD-B-L2-2026-11-20.md` (ODL-01…ODL-23).

### 11.2 What the selection does not change here

**OBE-41 (Observed).** §3–§7 stand as written: the per-option obligations (OBE-10, OBE-14, OBE-18), the
silent-failure modes (OBE-11, OBE-15, OBE-19), the comparison (OBE-21) and the no-dominance finding (OBE-23).
None of them was a recommendation, and the choice retracts none of them.

**OBE-42 (Calculated — selected risk preference vs comparative finding).** These are three different kinds of
statement, and the record keeps them apart so that the evaluation cannot appear to have decided on evidence:

| Statement | Kind | Status |
|---|---|---|
| No option passes every sub-criterion better than the others; L2 wins S2/S4/S5, L1 wins S3 (OBE-23) | **comparative finding** | unchanged |
| An unwarranted state is least problematic, a truncated history that reads as complete is next, a write that confers status is least acceptable (ODL-18) | **owner's risk judgement** | **selected**, not derived here |
| This project adopts L2 | **decision** | **made 2026-11-20** (OBE-40) |

This record could pose the alternatives (OBE-29, refined by OBE-37); it could not produce the middle row, because
that row is a judgement about what this system can live with, not a fact about the options.

**OBE-43 (Calculated — the ordering maps onto the recorded modes).** The preference follows the surfacing axis of
OBE-22 exactly: L1's mode appears at **audit** (least problematic), L2's at **read** (middle), L3's at **write**
(least acceptable). It is a ranking of **tolerability**, consistent with CEV-40's point that likelihood and
recoverability are unequal — not a claim about which mode is rarest.

**OBE-44 (Inferred — the conditional is inherited, not discharged).** OBE-26/OBE-27 held L2's viability
contingent on distinguishing "no successor" from "unknown", and CEV-18 found that a closure declaration is
**required in effect by I15/I18 and absent from the text**. Selecting L2 converts that contingency into an
**inherited obligation**, which is a stronger requirement than a preference and weaker than a solution. OBE-14's
two obligations (retain the reason; close the stream) survive intact, and OBE-36's open point — whether a stream can
say "no further entries by design" — remains **Unestablished**.

**OBE-45 (Observed — one tension recorded, not resolved).** L2's accepted residual failure is a silent, confident
default, which is the shape K-R17 rejects at write in the acquisition stream — *"a silent no-op hides formation
bugs"* (CEV-02/CEV-03). The streams differ in semantics, and CEV-05 shows the acquisition rule needed two parts.
Whether the lifecycle stream should carry a comparable rejection-at-write discipline is **not established**, and
answering it would be closure-mechanism design, which no record authorises.

**OBE-46 (Observed — boundaries carried, not restated as requirements).** The owner's accompanying boundaries — no
careless reconstruction, no relaxation of I15, no inference from absence, no fabricated evidence — are recorded at
ODL-16/ODL-20/ODL-23 and are not new requirements here (CEV-38).

**OBE-47 (Observed — everything else is untouched).** C′, OD-g, ND-1, DN-4, §8.1, the `active`/`established`
vocabulary question and the prohibition on designing Howl's state machine remain exactly as OBE-24, OBE-32,
OBE-33 and OBE-34 state. The preference implies no field, event type, relation or mechanism.

**End of record.**
