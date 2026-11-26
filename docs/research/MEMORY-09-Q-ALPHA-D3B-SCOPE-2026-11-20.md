# MEMORY-09 — Q-α: Does D-3b Reach Lifecycle Transitions?

## A bounded examination of the decision's wording and decision scope

**Status:** scope examination. **Q-α is not answered here and is not decided here.** No schema, carrier, field,
event type, descriptor, predicate, marker, mechanism or implementation. **No new requirement, no new authority, no
retention commitment.** Read-only: no canonical edit, no test, no probe, no ledger entry. **ND-1**, the **validity
predicate**, **Q-β** (ODC-18) and **ODL-39** remain OPEN. All three retention positions remain available.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | **Q-α** — does D-3b's cause-retention obligation apply to **lifecycle transitions**? |
| Method | the decision's **wording** and **decision scope** only — *not* inferred from the fact that lifecycle events can be statement-less |
| Answers | what class the decision names (§2), what the decision was about (§3), the two available readings of its class (§4), and the status of Q-α as a question (§5) |
| Decides | **nothing.** Q-α is left where the text leaves it; no position selected, ranked, merged or dropped |
| Does not touch | **Q-β** (ODC-18), ND-1, the validity predicate, D-3b's text, the 118 requirements, P-G0/P-G1/P-G2 |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Narrowed** · **Superseded**.

---

## 1. Why this is a scope question and not a retention decision

**QAX-01 (Observed — the question, as posed).** Whether **D-3b's** cause-retention obligation reaches **lifecycle
transitions**. It was left candidate-level by TGP-11 and listed as unestablished by GRP-51, and PGR-25 made it one of
the two questions P-G0's coherence now rests on.

**QAX-02 (Observed — the method the owner set, and why it matters).** The examination is to use **the wording and
decision scope of D-3b**, and explicitly **not** to infer applicability "merely because lifecycle events can be
statement-less." That exclusion is doing real work: statementlessness is a property lifecycle events *can* have
(TGP-11; CLV-41), so a property-based inference would decide Q-α in one step. It is excluded because **a property an
event can have is not evidence about what a decision was made about.**

**QAX-03 (Calculated — what would count as an answer).** Q-α is answered only if D-3b's text or decision scope
supplies one of: an explicit statement of reach; a class definition that is stream-specific; or a scope term in the
decision that excludes or includes lifecycle. Absent all three, the honest answer is *not established* — and that
answer is not a null result, because it changes what kind of question Q-α is (§5).

**QAX-04 (Observed — what this record does not do).** It does not decide Q-α, does not widen or narrow D-3b, does
not create a requirement, and does not lean on Q-β. Where the text is silent, silence is recorded as silence.

---

## 2. What the decision's words name

**QAX-05 (Source-derived — the decision, verbatim).**

> "MEMORY-09 may create events without a source statement. Such events must retain their actual cause and
> independently identifiable occurrence. Absence of conversational evidence must never be replaced with fabricated
> evidence or an inferred statement position."

**QAX-06 (Calculated — three sentences, and the class is stated once).** The first is a **permission**; the second is
the **obligation**, attached to "such events"; the third is a **prohibition**. The class appears exactly once, in the
first sentence, and every later sentence inherits it by anaphora. There is no second, narrower statement of the class
anywhere in the decision.

**QAX-07 (Calculated — the class is a property of an event, not a name for a stream).** "Events without a source
statement" identifies its members by something an event **lacks**. It names no stream, no record type and no
operation. Structurally, then, the class is **open**: any event that lacks a source statement satisfies it, in
whatever stream such an event occurs.

**QAX-08 (Observed — the record's own vocabulary).** A case-insensitive search of the decision record finds
**zero occurrences of "lifecycle"**. The substantive entries use acquisition and event-schema vocabulary:
`evidence_span`, `via`, `(turn, span position)`, "the event schema", "acquired". The words *transition*,
*supersession*, *retraction*, *withdrawal*, *promotion* and *status* appear only in the label legend and in one
table row that holds the status fold **separate** (QAX-14).

**QAX-09 (Calculated — what the wording does and does not settle).** It settles that the class is not defined by
stream (QAX-07). It does **not** settle that the class was *intended* to span streams — those are different claims,
and only the first is textual. The gap between them is the whole of §3.

**QAX-10 (Observed — the anaphor's reach).** "Such events must retain their actual cause" is a statement about the
members of a class, not about a stream's records. So the obligation's **subject matter** is event-level. Again: this
is a fact about grammar, not about scope.

---

## 3. What the decision was actually about

**QAX-11 (Source-derived — the problem it addresses).** SD-07: the decision "closes a gap between what T2 already
permits (`narrator_grant`, ambient knowledge) and what **the event schema assumed (that every event points at a
span)**". The assumption it removes is a **schema** assumption, and the schema named is the acquisition event's
(SD-06/07).

**QAX-12 (Source-derived — the two motivating provisions).** SD-06 cites `via` including `narrator_grant` — "a grant
is an act, not a character's utterance" — and S32's ambient availability, a newcomer receiving knowledge "without
acquisition". Both are **knowledge-acquisition** provisions. Neither is a lifecycle provision, and no lifecycle
provision is cited anywhere in the decision.

**QAX-13 (Source-derived — the consequences it draws, all acquisition-side).** SD-08: `(turn, span position)` is
undefined for such events. SD-10: "sort it at the end of the turn" is forbidden as an inferred statement position.
SD-12: a total order needs a key every event carries. SD-18: "`via` names the *kind* of acquisition; it does not
name the *actor or mechanism* that caused the event." SD-19: `evidence_span` must be nullable. Every consequence is
stated over **acquisition-log carriers** — span, turn, `via`, the evidence field, the reconstruction fold (D-4).

**QAX-14 (Source-derived — the effects table treats the status fold as separate).** §6 lists the decision's effects
on D-1, D-2, D-3a, D-3b, D-4, then records **"ND-1 (status fold) | unchanged, still separate"** and DN-6a/DN-6b
"**still open**" with the note that the decision says events must record their cause, it does "**not** say who may
cause them". The record's own accounting places the status fold outside what was decided, and marks the authority
question as untouched.

**QAX-15 (Source-derived — the canonical surface it would have touched).** §9 proposes **no** edit to r33 §6.2,
records that the 118 requirements are unchanged, and files SD-02's cause requirement as an obligation **without a
conformance row** (SD-24). A decision that had reached lifecycle transitions would have had to say something about
K-R12's stream; it says nothing, and §9 proposes nothing there.

**QAX-16 (Observed — the bounded check on the requirements document).** A keyword check of
`docs/architecture/REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md` returns **no** occurrence of "D-3b" or
"statement-less". No lifecycle requirement states this obligation. That is consistent with SD-24 and adds no reach.

**QAX-17 (Calculated — what §3 establishes).** The decision's **scope** — its motivation, its consequences, its
effects table and its proposed canonical changes — is entirely **acquisition-side**. Nothing in the record was
written about lifecycle transitions, and no lifecycle provision was consulted in reaching it.

**QAX-18 (Calculated — and what §3 cannot establish).** None of that is a statement of **exclusion**. The decision
never says "and not lifecycle transitions". Its class, as written, is stream-neutral (QAX-07). So §3 shows what the
decision **was about**; it does not show what its words **exclude**. A reader may not turn an absence of discussion
into a boundary — that is the error I15 names, applied to scope.

---

## 4. Two readings of the class, stated fairly

**QAX-19 (Calculated — R-prop: the class is property-based, so occurrence is enough).** If the class is defined
solely by lacking a source statement (QAX-07), a statement-less lifecycle transition **is** a member, and D-3b's
obligation attaches to it without further decision. Supporting this reading: the class is stated once and never
restricted; "independently identifiable occurrence" (SD-03) is stream-neutral; the decision's subject is "events" in
an architecture that has more than one event stream; and the prohibition (SD-04) is phrased about *evidence*, which
is equally relevant wherever an event records none. **Under R-prop, applying D-3b to a statement-less lifecycle
transition is application, not extension** — and P-G1 is a description of an existing obligation, not a request for
new authority.

**QAX-20 (Calculated — R-stream: the class is the events the decision was addressing).** If a class's membership is
fixed by the decision's scope as well as its words, the class is the **acquisition events** that lack a span, and
lifecycle transitions are outside it. Supporting this reading: every motivating provision, every consequence and
every proposed canonical effect is acquisition-side (QAX-11…QAX-15); the record's own effects table holds the status
fold separate (QAX-14); and the ledger entry that would have recorded the decision was never made (§9), so its scope
was never fixed in canonical terms. **Under R-stream, applying D-3b to a lifecycle transition is an extension** by
the reader — new authority, not existing.

**QAX-21 (Calculated — the split, stated precisely).** The **wording** is property-based and stream-neutral
(R-prop); the **decision scope** is acquisition-specific (R-stream). Neither reading is excluded by the other's
evidence, because they answer different questions: *what do the words admit?* and *what did the decision govern?*

**QAX-22 (Observed — which reading the owner's exclusion forbids choosing).** Preferring R-prop **because** lifecycle
events can be statement-less is exactly the inference excluded in QAX-02. The property's applicability to lifecycle
events is what Q-α asks about; it cannot also be its answer. So R-prop is available only as a **reading of the
wording**, never as an inference from the property.

**QAX-23 (Observed — and R-stream is not established either).** Nothing states exclusion (QAX-18). Choosing R-stream
because the record happens not to discuss lifecycle transitions would be the same error mirrored — reading absence of
discussion as a boundary.

---

## 5. The status of Q-α

**QAX-24 (Observed — the answer to the question as posed).** **Q-α is not answered by D-3b's wording or decision
scope.** The decision states neither inclusion nor exclusion, and no provision elsewhere in the corpus was found
that states either. TGP-11's "no text states which lifecycle transitions fall inside it" stands, now with the reason
attached: the text's class is property-based while its scope is acquisition-side.

**QAX-25 (Calculated — what kind of question Q-α therefore is).** It is not a question with a discoverable answer
in an existing document. It is an **owner scope decision**: whether an obligation fixed for one event stream is to
be **extended** to another. Under R-prop the extension is unnecessary; under R-stream it is the only route. Either
way the answer is not derivable, and the record that would fix it (O-16) was never made.

**QAX-26 (Calculated — the "authority" half, stated exactly).** The scope-and-authority question decomposes into
two parts that are often run together:

| Part | Question | State |
|---|---|---|
| **Scope** | did the decision purport to govern lifecycle transitions? | **no statement either way** (QAX-17/QAX-18) |
| **Authority** | if it did not, may a reader apply it there? | **no** — a reader cannot extend a decision's reach; that is the owner's (ODL-31's obligations-not-authorisations logic; ODC-31's discipline) |

So the authority half is determinate **only conditionally**: *if* R-stream holds, application is unauthorised
extension; *if* R-prop holds, no extension is needed. The condition is the scope half, and the scope half is
unstated.

**QAX-27 (Observed — Q-α is not ND-1, and must not be read as it).** ND-1 asks **who may cause** a lifecycle change
— an authority condition on actors, still OPEN (ODL-32; ODB-57). Q-α asks whether a **retention obligation** applies
to a class of events. D-3b itself keeps these apart: SD-05 and SD-23 state that recording a cause "presupposes that
*something* caused it, not that the something is authorised", and §6 records DN-6a/DN-6b as still open. Answering
Q-α would not answer ND-1, and ND-1's openness does not answer Q-α.

**QAX-28 (Observed — no position is implied by this finding).** That Q-α is an owner scope decision does not make
P-G0 correct, does not make P-G1 an extension, and does not make P-G2 necessary. It identifies **which decision**
the owner would be taking if Q-α were taken **— and this record does not recommend taking it, or recommend an
answer.**

---

## 6. What this changes for the distinctness question

**QAX-29 (Calculated — P-G1's name already presupposes Q-α's answer).** P-G1 is called "**D-3b's floor**". Under
R-prop that name is accurate and P-G1 adds **nothing** — it describes an obligation that already binds. Under
R-stream the name is a **proposal**, and P-G1 would be the extension itself, not a description of existing
requirements. So the owner's observation in the prior exchange — that P-G1 "isn't necessarily asking the owner to
authorize more retention" and "may simply describe the existing cause-retention obligation more accurately" — is
**exactly conditional on Q-α**, and not on anything else. That is the sharpest reason Q-α comes first: **P-G1's
status as description-or-proposal is not determinable until it is resolved.**

**QAX-30 (Calculated — the two branches for the option space, unranked).**

| If Q-α resolves… | P-G0 | P-G1 | Positions live |
|---|---|---|---|
| **R-prop** — D-3b already spans events | would require **not applying** an existing obligation to lifecycle transitions — not a neutral silence | **descriptive**; adds no authority | P-G0 becomes a claim about non-application; P-G1 describes |
| **R-stream** — D-3b is acquisition-only | **live and unburdened**; displaces nothing | **the extension proposal** itself | both live as genuine alternatives |

Neither branch is chosen. Recording them shows that Q-α does not decide the retention question — it decides **what
the existing positions mean**, which is what the reassessment asked.

**QAX-31 (Observed — Q-β is untouched and independent).** Whether an actor-reported rationale is a ground (ODC-18)
is unaffected by Q-α: it is a classification of content, and it arises under D-3b wherever D-3b applies. Nothing here
leans on it, and Q-α's resolution would not settle it. The two must not be collapsed into one owner decision
(PGR-25/PGR-26 stand).

**QAX-32 (Observed — no retention commitment is made).** This record adds no obligation to retain anything, creates
no field, widens no requirement, and does not convert Q-α into a task with an authorised outcome. If the owner later
takes Q-α, that decision would need its own record; this one does not pre-write it.

**QAX-33 (Observed — all three positions retained).** P-G0, P-G1 and P-G2 remain available, unselected, unranked and
undropped (PGR-24/PGR-28). Axis H remains identical across them (GRP-31); the validity predicate and ND-1 remain
open; ODL-39 and lifecycle completeness are untouched.

---

## 7. Traceability

**QAX-34 (Observed).**

| What | Where |
|---|---|
| D-3b verbatim, and its four settled parts | SD-01…SD-04; TGP-09 |
| That recording a cause is not an authority decision | SD-05, SD-23; §6's DN-6a/DN-6b row |
| The decision's motivating provisions and schema assumption | SD-06, SD-07 |
| Its acquisition-side consequences | SD-08…SD-13, SD-17…SD-19 |
| Its effects table, including the status fold held separate | §6 (D-1…D-4, ND-1, DN-6a/DN-6b) |
| No canonical change proposed; no conformance row | §9; SD-24 |
| Reach as candidate-level, and per-event not per-stream | TGP-10, TGP-11 |
| Q-1's three-valued answer | TGP-12; TGP-13 |
| P-G0's availability resting on Q-α | PGR-10, PGR-25, PGR-26 |
| Q-α/Q-β independence | PGR-25/PGR-26; PGR-27 |
| Q-β as the residual classification | ODC-18, ODC-19; CGE-43 |
| ND-1 as the actor-authority question | ODL-32; ODB-57 |
| Obligations are not authorisations | ODL-31 |

---

## 8. Limits and record status

**QAX-35 (Observed — method).** **Newly read for this record:** the D-3b decision record in full (206 lines), TGP-08
through TGP-13 verbatim, and a keyword check of the lifecycle requirements document. **Not inspected:** r32's bytes,
the frozen preregistration, v9 beyond §7, and any implementation's cause handling. No test, probe, benchmark or model
call was run (QAX-36).

**QAX-36 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The question is textual and decision-scoped; it needs no execution
and none was attempted.

**QAX-37 (Observed — the labels' limits).** QAX-06…QAX-09 and QAX-17…QAX-22 are **Calculated** — they follow from the
decision's own text and structure, not from a new source. QAX-24 is **Observed** as an absence and is the record's
central claim. QAX-26's authority half is **conditional on purpose**. Nothing here is presented as a canonical
finding in r33 or v9: no provision states D-3b's stream reach, which is the finding.

**QAX-38 (Observed — the search's limits).** The vocabulary check is lexical. It could miss a scope statement phrased
without the words searched, so it was supplemented by reading the decision in full and its effects table (QAX-14) —
which is the place a stream reach would have been recorded had it been intended. The requirements-document check
(QAX-16) is likewise lexical and is reported only as corroboration of SD-24.

**QAX-39 (Observed — no other conflict found).** The decision's text, its motivation, its consequences and its
effects table are internally consistent and consistent with TGP-10/TGP-11's reading. The seam this record records is
QAX-21: **the words are stream-neutral, the decision was acquisition-side**, and nothing reconciles them in the
corpus.

**QAX-40 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Obligations created or widened | **0** |
| Q-α, Q-β, ND-1, validity predicate, ODL-39 | **all OPEN, un-narrowed** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; no canonical capture, no ledger entry |

**End of record.** **Q-α is not answered.** D-3b's **wording** names a stream-neutral class ("events without a source
statement") while its **decision scope** is entirely acquisition-side; the corpus states neither inclusion nor
exclusion. Q-α is therefore an **owner scope decision** — whether to extend an existing cause obligation to lifecycle
transitions — and not a fact recoverable from the decision (QAX-24/QAX-25). **P-G1's description-or-proposal status
turns on it** (QAX-29). **Nothing is decided, no position is selected or dropped, Q-β and ND-1 remain open, and no
retention commitment is made.**
