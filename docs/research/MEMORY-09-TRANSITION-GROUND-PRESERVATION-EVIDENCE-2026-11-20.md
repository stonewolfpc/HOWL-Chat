# MEMORY-09 — Transition-Ground Preservation: What Carrying a Reason Can and Cannot Establish

## Bounded evidence on ground retention under L2, the authority boundary, and I17's actual scope

**Status:** evidence record. **No schema, carrier, field, event type, mechanism, predicate or implementation.** No
new requirement, no new authority. Read-only: no canonical edit, no test, no probe, no ledger entry. OD-b remains
**settled as L2** (ODL-01). ODL-26 (closure), ODL-28 (determinism), ODL-29 (`active`/`established`) and ODL-30
(§8.1) are **untouched**; ODL-39 remains **OPEN**.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Scope | Q-1…Q-6 as bounded by the owner: whether L2 requires ground retention; what carrying a ground means; what I17 prohibits; whether a projection can preserve a reason without reading it as proof; what CLV-27/28 establish; implications for uncertainty and auditability |
| Excluded | schema, carrier, field, event type, mechanism, closure design, implementation, any new authority inferred from the L2 decision |
| Basis | r33: §3 (K-R2, K-R12, K-R13), §4 (I7, I15, I17), §10.2.1, §11, §11.3 `[P]`; D-3b record (SD-01…SD-19); ODL-27; CLV-27/28/30/41; CEV-21…CEV-24; LCM-03/LCM-31 |
| Prior records | `...OWNER-DECISION-OD-B-L2...` (ODL), `...SECTION-5-11-PROVENANCE...` (CLV), `...OWNER-DECISION-D3B...` (SD), `...CLOSURE-AND-VALIDITY-EVIDENCE...` (CEV), `...LIFECYCLE-COMPLETENESS-EVIDENCE...` (LCM) |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

**Notation:** `[X]` = the r33 provision carries that tag; `[P]` = the provision is in a section r33 itself marks
*proposed, not adopted*. Nothing in a `[P]` section is cited below as existing policy.

---

## 1. Scope, and why "ground" has to be separated before it can be answered

**TGP-01 (Observed — what this record answers).** Six questions, kept separate as the owner framed them. The
boundary under investigation is one sentence: **preserving why a transition was recorded must not silently become a
justification for treating that transition as authorised.** The record's job is to locate that boundary in canon and
in the prior records, not to build anything across it.

**TGP-02 (Inferred — five things currently called "ground").** The L2 discussion uses *ground*, *reason*, *cause*,
*basis* and *warrant* as though they named one thing. They are five distinct claims, and the whole question turns on
which is which:

| Term | The question it answers | Status in the corpus |
|---|---|---|
| **Cause** | what produced this event? | owner-decided for statement-less events (D-3b, SD-02); shape unsettled (SD-18) |
| **Ground / reason** | what basis was claimed for this change? | **not a canon term**; no requirement (TGP-08) |
| **Evidence** | what in the conversation supports the record? | `evidence_span`; declared non-nullable, defined in no schema (D-1) |
| **Warrant** | was the causing act entitled to change the state? | **ND-1**, open and widened |
| **Derivation basis** | what did the fold actually use? | constrained by K-R2 and ODL-28; not specified |

**TGP-03 (Calculated — the distinction that does the work).** A **ground** is a claim *about* warrant. A **warrant**
is the thing itself. No amount of ground becomes a warrant by being recorded, because the ground is asserted by the
party whose act it explains, and the warrant is established under conditions that party does not control. I17 states
exactly that asymmetry (TGP-25/26).

**TGP-04 (Inferred — evidence is not a ground either).** D-3b's own table separates them: evidence answers *"what in
the conversation supports this?"* and cause answers *"what produced this event?"* (SD-17). Neither answers *"why was
this change warranted?"* A span shows what a statement said; it does not show that the statement met the promotion
condition for the kind. So evidence cannot be substituted for a ground, and their presence is independent.

**TGP-05 (Observed — "ground" is not canon vocabulary).** A term search of r33 returns **no occurrence of
"ground"**. The word "reason" occurs once, at `:554`, inside §10.2.1 — a `[P]` section — as a field label in a
proposed evaluation record (`verdict + reason code`); the only other match is the unrelated word "reasoning" in
§13.1 (`:1009`). There is therefore no canonical answer to "what does r33 require of a transition's ground",
because r33 does not use the concept.

**TGP-06 (Inferred — why that is a finding, not a gap in the search).** A requirement can exist without its
vocabulary, but here the adjacent provisions are explicit and would have used it: K-R12 enumerates what the stream
records (`:111`), K-R13 defines the one read-time join (`:112`), and §11 specifies promotion conditions per kind
(`:600–612`). All three state *what* changes and *when*; none states that the basis must be kept. The absence is
consistent across the text, not an artifact of one keyword (TGP-53).

---

## 2. Q-1 — does L2 require transition grounds to be retained?

**TGP-07 (Source-derived — K-R12, `[X]`, r33 `:111`).**

> "Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus
> `scope` and `knowledge_class` changes — and a memory's current status, class and scope are a **projection of that
> stream**, not mutable authoritative fields."

The stream is defined by the **kinds of change** it carries. The clause constrains *where* state lives (projection,
not field). It says nothing about the content of an entry beyond the change it records.

**TGP-08 (Observed — the stream-wide answer from canon).** **Nothing in r33 requires a lifecycle transition to
retain why it occurred.** K-R2 (`:93`) requires reconstructability *from the events that exist*; K-R12 requires the
kinds; §11 (`:600–612`) requires the conditions *under which promotion happens*, which is a rule the system applies,
not a fact the entry must carry. L2 does not create the requirement (consistent with ODL-27), and canon does not
supply it.

**TGP-09 (Source-derived — an existing owner decision that predates OD-b).** D-3b (2026-11-19), verbatim:

> "MEMORY-09 may create events without a source statement. Such events must retain their actual cause and
> independently identifiable occurrence. Absence of conversational evidence must never be replaced with fabricated
> evidence or an inferred statement position."

**SD-02** reads that as: *such an event must retain **its actual cause***. That is a cause requirement, decided
before L2 was selected, and independent of it.

**TGP-10 (Inferred — its reach into the lifecycle stream).** CLV-41 records that **every** lifecycle transition —
entry, supersession, correction, retraction, forgetting — is a *candidate* statement-less event. If D-3b reaches
them, then for that subset the answer to Q-1 is not "the owner has not decided": cause retention is **already
required** by an earlier decision, and ODL-27's "L2 may decline to carry the ground" is unavailable there.

**TGP-11 (Unestablished — and the scope is per-event, not per-stream).** Whether any given transition is "without a
source statement" is a property of the event, not of the stream. A retraction caused by a character's correction has
a statement; a retraction caused by an out-of-band review finding does not; a supersession caused by a new fact may
have a statement, and one caused by expiry or by a narrator act may not. So D-3b's requirement is a **floor that
applies unevenly across the stream**, and no text states which lifecycle transitions fall inside it. That question
is not answered here and is not answered anywhere in the corpus.

**TGP-12 (Calculated — Q-1's answer is three-valued, not yes/no).**

| Class of transition | Ground retained? | Required by |
|---|---|---|
| **Statement-less** | **cause** — yes, per D-3b | owner decision, SD-02 |
| **Statement-bearing** | **not established** — no text requires any basis | none; evidence does not substitute (TGP-04) |
| **Stream-wide (all transitions)** | **no requirement** | K-R12 is silent (TGP-07/08) |

**TGP-13 (Observed — a qualification to ODL-27, not a correction of the decision).** ODL-27 states that "L2 may
decline to carry the ground, at the cost that 'why' becomes unrecoverable — that is a decision still to be made".
That remains accurate for the stream-wide question and for statement-bearing transitions. It is **not** accurate for
statement-less transitions if D-3b reaches them (TGP-10/11). The decision record did not read D-3b's cause clause
against the lifecycle stream; this record adds that reading and nothing else. **ODL-27's decision — that the
stream-wide question is the owner's — is unchanged.**

**TGP-14 (Unestablished — whether D-3b's "cause" is a ground at all).** SD-18 says a cause descriptor "names the
*actor or mechanism* that caused the event" and that its shape is unsettled. A cause and a ground can coincide (a
review finding causes a retraction *and* is its basis) and can diverge (a new fact causes a supersession; the basis
is that the old record no longer applies). **D-3b may therefore deliver cause without delivering ground**, and the
two must not be treated as interchangeable when the owner's question is about "why".

**TGP-15 (Inferred — the honest statement of Q-1).** L2 requires no ground retention. No canon provision requires
it. One prior owner decision may require **cause** retention on a subset whose boundary is undrawn. Whether the
stream as a whole carries grounds therefore **remains an owner-level decision**, exactly as ODL-27 and ODL-31 say
— with the qualification that the decision's floor is not zero (TGP-12).

---

## 3. Q-2 — what does carrying a ground mean, and what would make it an authority source?

**TGP-16 (Inferred — carrying, stated without a carrier).** To carry a ground is for a transition to retain a
resolvable reference to a fact stating **why the transition occurred**, such that a later reader recovers the stated
basis without re-deriving it from present state. Three things it is not:

- **not the transition kind** — `superseded` says *what* changed, not why;
- **not the evidence span** — that says what supports the record, not why the change was made (TGP-04);
- **not the derived state** — a projection output is a conclusion, not a basis.

**TGP-17 (Calculated — the two questions any carry must answer).** These are analytic axes, not options to select:

| Axis | The question | Why it matters |
|---|---|---|
| **Locus** | is the ground *inside* the transition, or *referenced by* it? | the closer the ground sits to what the fold reads, the smaller the step from "read" to "rely on" (CLV-28) |
| **Interpretation** | is the ground carried *as stated*, or *evaluated* by the fold? | interpretation is the exact point where a basis becomes a warrant (TGP-18b) |

CLV-28 already places the constraint on the first axis ("referenced by … never read as"). This record does not
choose a locus or an interpretation; it records that a representation which answers neither question cannot be
assessed against I17 at all.

**TGP-18 (Inferred — the three ways a ground becomes an authority source).** Precisely three, and they are
distinguishable:

1. **By presence.** The fold treats "a ground exists" as sufficient warrant — so the ability to append a reason
   becomes the ability to change status. That is the inversion I17 names, arriving through the annotation.
2. **By content.** The fold computes the state or the admissibility from the ground's kind or text — "basis was a
   correction, therefore `corrected`". The reason now authorises the status.
3. **By substitution.** The ground stands in for the validity predicate that does not exist. The stream is asked
   *"was this act entitled to change the state?"* and answers with *"here is why the act was done"* — a different
   question.

**TGP-19 (Calculated — substitution is the live one, and why).** ND-1 is open: "who may cause a lifecycle change is
not established" (ODB-57; ODL-32). Where no authority predicate exists, **any** recorded ground is the only
candidate answer to "was this warranted?", so the substitution is not a mistake a careless implementer must make —
it is the default that fills a vacuum. This is the mechanism behind CLV-27's warning that the inversion is "live
for L2, not only L3", and it is why the constraint cannot be left implicit.

**TGP-20 (Inferred — the structural reason a ground is weaker than it looks).** A ground is **self-reported**: it is
authored by the process or party whose act it explains. I17's requirement is that authority trace to "the conditions
under which it was established" — conditions the establishing party does not get to certify by describing them. So a
carried ground can raise **auditability** (what was claimed, by whom, in what form) but cannot by itself raise
**reliability** (that the claim held). Conflating those two is the error the owner's boundary is drawn against.

**TGP-21 (Source-derived, `[P]` — the closest existing expression, in a different domain).** r33 contains the idea
that a system cannot certify its own completeness: K-R83 (`:778`, §11.3, **proposed, not adopted**) — "a registry
cannot prove its own completeness … a missing entry is not evidence that nothing changed" — and K-R84 (`:779`,
`[P]`), which requires every subsystem to state the independent evidence that would expose its own omission. These
are the same shape as TGP-20 and they are **not adopted**, and their scope is identity provenance and omission
detectors, not lifecycle grounds. They are cited as evidence that the architecture can conceive of the property, and
for nothing else.

**TGP-22 (Observed — no carrier is implied by any of this).** Whether a ground is a field, a relation, a descriptor,
part of an event envelope, or a separate record is **not decided here** and is not decided in any prior record.
ODL-27 leaves the representation open; SD-18 leaves the cause shape open; SD-05 explicitly declines to settle
"whether a cause is a field, a descriptor or a relation".

**TGP-23 (Unestablished — a ground without a provenance for the ground).** Nothing in the corpus says whether a
carried ground must itself record **who stated it** and **under what conditions**. Without that, a ground is
indistinguishable from the system's own inference — which is exactly the substitution D-3b's SD-04 forbids in the
evidence dimension ("an inferred statement position"). The analogue for grounds has not been decided, and it is the
sharpest open point in this section.

**TGP-24 (Inferred — the boundary, stated once).** Carrying a ground becomes an authority source at the moment
**anything about the derived state or the transition's admissibility depends on the ground**. Before that moment it
is provenance; after it, it is a warrant the stream has granted itself. The boundary is therefore not in the
ground's content but in the fold's dependence on it.

---

## 4. Q-3 — what does I17 actually prohibit?

**TGP-25 (Source-derived — I17, `[X]`, r33 `:305`).**

> "**A result earns authority from the conditions under which it was established, not from when or where it was
> written.** Applies beyond temporal comparison to retrieval, formation and conflict evaluation: never 'it was
> computed most recently', always 'it was computed under these conditions'."

**TGP-26 (Calculated — four things it prohibits).**

1. **Authority from the act of writing.** "Where it was written" is the record's location; recording is not
   establishing.
2. **Authority from recency.** "Computed most recently" is explicitly the wrong test.
3. **Authority from position.** Order in a log is a fact about the artifact, not about the world.
4. **Generally: substituting any fact about the record for a fact about the conditions.** Clauses 1–3 are instances
   of this one.

**TGP-27 (Calculated — what it requires, and its limit).** Authority must be **traced to the conditions of
establishment**. I17 does not enumerate those conditions, name who evaluates them, or say what happens when they are
unknown — CEV-24's finding, which stands: "the principle is present; the **predicate** is not."

**TGP-28 (Inferred — applied to grounds, the prohibition is exact).** "The transition records a reason, therefore
the transition was warranted" is prohibited, because the recording is a fact about the record and the warrant is a
fact about the conditions. This is CEV-23's reading of I17 applied to the ground dimension rather than to the act of
recording an act: the same error, one carrier further out.

**TGP-29 (Source-derived — I7, `[X]`, r33 `:293`).** "Retrievability does not confer authority." A ground that is
stored and can be retrieved adds no standing by that fact. I7 and I17 are siblings here: one blocks retrieval as
authority, the other blocks recording as authority.

**TGP-30 (Calculated — what I17 does *not* prohibit, stated to prevent over-reading).**

| Not prohibited | Why |
|---|---|
| Recording a reason | I17 constrains derivation from the record, not what the record may contain |
| A transition and its ground being stored together | that is a representation choice I17 does not reach |
| Different grounds producing different authority **when the grounds are the conditions** | I17 requires that link; it does not forbid it |
| Recording an unauthorised transition along with its reason | recording a bad act is not the error; treating the record as its warrant is |

**TGP-31 (Inferred — therefore I17's force lands on the fold).** A writer may record a ground freely. The
constraint bites at one point only: **the derivation must not read it as standing**. That is why the separation is
statable as a constraint before any design exists, and why it does not require a mechanism to be meaningful.

---

## 5. Q-4 — can the projection preserve a transition's reason without treating it as proof of authority?

**TGP-32 (Calculated — yes, and the separation is a property, not a design).** The requirement is expressible
without choosing a carrier: **the derived state must be invariant under changes to grounds alone.** Everything the
fold emits — the state name and its uncertainty — depends on the transitions; the grounds travel alongside and
affect neither.

**TGP-33 (Calculated — the discrimination test, a criterion and not a mechanism).** Two streams with **identical
transition sequences** and **different grounds** must fold to the **same** state. If they do not, the reason has
become an authority source. The test needs no field, no marker and no implementation:

| Perturbation | Required result |
|---|---|
| grounds differ, transitions identical | identical derived state |
| a ground is removed, transitions identical | identical derived state (not a rejected transition) |
| a ground is added, transitions identical | identical derived state |

**TGP-34 (Inferred — the test's limit, which must be stated with it).** It detects **state-level** leakage only. A
fold could read a ground in a way that never changes the emitted state — annotate confidence, suppress a warning,
choose between two equal-state paths — and pass all three rows. So the test is **necessary, not sufficient**. It is
still the strongest available check, because it is the only one that can be applied without a representation.

**TGP-35 (Inferred — and it constrains the representation without selecting one).** The test is only usable if
grounds are **separable** from transitions — a stream in which the ground is fused into the transition kind cannot
vary them independently, so it cannot demonstrate that it does not read them. That is a property the representation
would have to have; naming the property is not naming the representation.

**TGP-36 (Inferred — preserving the reason requires preserving disagreement).** If the ground says "replaced by X"
and the transition kind is `corrected`, a projection that silently prefers one has destroyed the evidence that
something was wrong. Preservation's real yield is not a tidier history; it is a history in which a **mismatch between
the stated basis and the recorded change** remains visible. That is also the only way a later auditor can find the
class of error the owner ranked most severe.

**TGP-37 (Inferred — a reason must be preserved as asserted, not as established).** The record must be able to
distinguish "a ground was stated" from "a ground was verified". Today **nothing** in r33 can express the second, for
any dimension — C′ is open and untouched by OD-b (CLV-32; ODB-48). So a preserved ground is, by construction, a
claim. Whether the representation labels it as such is a labelling property, not a field, and it is not decided here.

**TGP-38 (Unestablished — whether ground-blindness and D-3b are jointly satisfiable in practice).** Formally they
are compatible: the cause is carried, the fold ignores it. But no text states the compatibility, no text says an
auditor may read what the fold ignores, and the point at which the auditor's reading becomes a second derivation is
undecided. This is where an owner decision would be required, and it is not taken.

---

## 6. Q-5 — what do CLV-27 and CLV-28 establish, and what stays unresolved?

**TGP-39 (Source-derived — CLV-27, quoted).**

> "**L2 does not by itself retain why a transition occurred.** The relation carries the type of change, not its
> justification. Retaining 'why' under L2 requires the transition to **reference** a ground, and that reference is a
> design choice — one that brings the cause into the transition record. The inversion the owner flagged is therefore
> **live for L2**, not only L3: a transition that *carries* its cause can drift into letting the cause authorise the
> status."

**TGP-40 (Source-derived — CLV-28, quoted).**

> "Cause must be **referenced by** the transition, never **read as** the status. A representation that cannot state
> a transition without also asserting its ground's authority has merged two of the six dimensions in CLV-30."

**TGP-41 (Calculated — four things the pair establishes).**

1. The retention gap is **real for L2**, not an L3 artifact. Selection of L2 does not close it.
2. The inversion risk attaches to the **reference**, not to the option — so it survives any option that carries a
   ground.
3. The separation is **statable as a constraint** (referenced-by / never-read-as) before any carrier exists.
4. Retention is a **design choice**, not a property of the stance. L2 neither requires nor forbids it.

**TGP-42 (Inferred — five things they leave unresolved).**

| Question | State |
|---|---|
| Is retention **required**? | Q-1: no stream-wide requirement; D-3b may require **cause** on a subset (TGP-12) |
| A reference to **what**? | unshaped — a fact, an event, an actor, an uninterpreted statement (SD-18; TGP-16) |
| Is the ground the **Cause** dimension or a new one? | CLV-30 lists Cause as "why"; if the ground is Cause, no seventh dimension is needed — nothing decides this |
| Can "never read as" be **enforced**, or only reviewed? | no mechanism is authorised, so today it is a **review obligation**, not an invariant |
| Does the constraint survive **closure**? | a ground does not help the fold know whether the tail was lost (LCM-31); ODL-26 untouched |

**TGP-43 (Observed — CLV-27/28 stand, with one reading note).** CLV-27's phrase "brings the cause into the
transition record" is the closest sentence in the corpus to describing a carrier. It must be read as a statement of
what retention **entails**, not as a design: it says a reference is needed, not what a reference is. Nothing in
CLV-27/28 is withdrawn, and this record adds no option to their comparison.

---

## 7. Q-6 — implications for uncertainty and auditability, without selecting a representation

**TGP-44 (Calculated — the auditability gain, stated exactly).** Without a ground, an auditor can establish **what**
changed and **when**, but not **what was claimed as its basis**. With one, the auditor can compare the claim against
the applicable rule — **if a rule exists**. ND-1 is open, so retention converts "unanswerable" into "decidable
later". It does not make anything decidable now, and it does not make the transition more likely to be correct.

**TGP-45 (Inferred — the auditability cost, which is the sharpest honest finding here).** A recorded ground is an
artifact that **looks like** justification. Later readers will read it as one even if the fold never does.
Preservation therefore **exports the I17 problem from the code to the reader**: the system can be ground-blind and
still produce a record that a human interprets as a warrant. That is not an argument against retention; it is the
reason the CLV-28 constraint cannot be left to implementation taste.

**TGP-46 (Inferred — two uncertainties that must not merge).**

| Situation | What it establishes | Governing rule |
|---|---|---|
| **No ground recorded** | not established — **never** "no reason existed" | I15 (`:301`) |
| **Ground recorded, not evaluated** | a claim was made; its truth is unverified | I15 + C′ open |

A representation that shows a ground for some transitions and nothing for others has created a **new absence**, and
I15 applies to it with full force. This is the closure record's finding one level down: absence-of-ground is not
evidence about the existence of a reason, exactly as absence-of-event is not evidence about what happened.

**TGP-47 (Inferred — the failure-class observation, labelled as analogy only).** The owner ranked "recording an act
that confers status without authority" as **least acceptable**. A ground read as a warrant is that failure arriving
through the annotation rather than through the act — same **class**, different carrier. **This is a structural
analogy, not an inheritance:** the ranking was about derived state, and nothing here extends it to grounds, and
nothing here borrows authority from the L2 decision to require anything. It is recorded because it explains why
TGP-24's boundary is substantive rather than hygienic.

**TGP-48 (Inferred — what retention does not fix).** Truncation (a lost transition's ground is lost with it;
LCM-31). Authority (ND-1). Evidence (D-1). Occurrence (OD-g). The correctness of the derived state, which depends on
an admissible-transition set that does not exist (LCM-24; `active`/`established` open). A carried ground is evidence
*about* the transition, never a repair of any of these.

**TGP-49 (Unestablished — the owner's actual choice, unranked).** Three positions are consistent with everything
above, and this record does not order them:

1. **No ground** — "why" becomes unrecoverable; the stream stays cheapest and I17 cannot be violated by a ground
   because none exists.
2. **D-3b's floor** — cause on statement-less transitions only; the stream's retention is uneven by event class, and
   the unevenness itself becomes a fact readers must understand.
3. **Ground on every transition** — full auditability, full reader-misreading risk, and the largest surface on which
   the CLV-28 constraint must hold.

**TGP-50 (Observed — the decision this record does not make).** Whether L2 retains transition grounds, and at which
of TGP-49's positions, **remains the owner's**. ODL-27 called it "a decision still to be made"; that remains true,
now with a floor (TGP-12) and a stated boundary (TGP-24) attached.

---

## 8. What this record does not do, and its limits

**TGP-51 (Observed).** No schema, carrier, field, event type, descriptor, predicate, marker, mechanism or
implementation is proposed. No position in TGP-49 is selected, ranked or preferred. Nothing here becomes a
requirement, and nothing here is a recommendation.

**TGP-52 (Observed — no borrowed authority).** L2's selection does not authorise ground retention, and this record
does not infer authority from it. ODL-26 (closure), ODL-28 (determinism), ODL-29 (`active`/`established`) and
ODL-30 (§8.1) are untouched. **ODL-39 remains OPEN**; this record adds nothing to the closure question and does not
treat K-R17 as inherited policy. D-3b is applied as written, not reopened or extended (TGP-13/14).

**TGP-53 (Observed — method limits).** The term search (TGP-05) is lexical and could miss a paraphrase; it was
supplemented by reading the adjacent provisions that would have carried the concept (TGP-06). r33's §11.3 rules
(K-R53–K-R87) are `[P]` — proposed, not adopted — and are cited only as such. Not inspected for a ground or cause
carrier: v9, r32's bytes, the frozen preregistration, or any implementation's event handling. D-3b's reach into the
lifecycle stream is inferred from CLV-41's candidate claim, not asserted by D-3b's text.

**TGP-54 (Observed).** No canonical file was edited; r33, v9, the owner ledger and the conformance map are
unchanged; the 118 requirements are unchanged; 0 tests were run.

**End of record.**
