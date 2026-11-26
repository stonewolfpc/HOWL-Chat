# MEMORY-09 — Q-β Decision Readiness

## When an actor reports why they caused a transition, is that rationale also a ground?

**Status:** decision-readiness record. **Research evidence only — not canonical policy, not a requirement, not a
selection.** Q-β is **OPEN**; this record does not answer it, does not lean toward either answer, and does not
recommend one. No schema, carrier, field, event type, descriptor, predicate, marker, mechanism or implementation.
Read-only: no canonical edit, no test, no probe, no ledger entry. Q-α and the cause classification remain settled.
**ND-1**, the **validity predicate** and **QAD-19** stay separate and untouched.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | **Q-β** — when an actor reports why they caused a transition, is that reported rationale **also a ground**, or does it remain a **cause accompanied by a report of the actor's rationale**? |
| Scope | the **report-only** case (`...CAUSE-GROUND-EQUIVALENCE...` CGE-24/CGE-25) — not every cause, and not every transition |
| Answers here | what each answer would mean (§3/§4), what each would change (§5), and what neither can change (§6) |
| Does not do | answer it, rank the answers, draft the decision, or create a requirement |
| Record status | **research evidence**; no canonical capture, no ledger entry, pending separate authorisation |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The question and its scope

**QBR-01 (Source-derived — the question, in the owner's words).** "When an actor reports why they caused a
transition, is that reported rationale also a ground, or does it remain a cause accompanied by a report of the
actor's rationale?"

**QBR-02 (Calculated — the question is correctly scoped, and the scoping does work).** It asks only about the case
where **an actor reports why**. That is the case the corpus calls **report-only**: "an operator's stated intention; a
narrator grant's purpose; a model's own decision rationale", where "the only access to the cause is the actor's
account" (CGE-24). The question does not reach causes that are independently observable.

**QBR-03 (Calculated — and for observable causes no question arises).** A clock expiry's cause "is not a claim about
why the change was warranted; it is a fact about elapsed time" (CGE-26). Nobody would call that a ground, and Q-β does
not ask about it. So the answer, whatever it is, is confined to the report-only case — which is why it cannot be
generalised to "all causes" without contradicting CGE-26.

**QBR-04 (Observed — report-only is not the same property as statement-less).** "A statement-less event may have an
independently observable cause, and a statement-bearing event's cause may be report-only" (CGE-29). The two boundaries
are distinct and **both undrawn**. So Q-β's extension is unestablished, and this record does not estimate it.

**QBR-05 (Calculated — what Q-β is not about).** Not the validity basis (conditions of admissibility — GRP-11), not
the warrant (ND-1), not evidence (`evidence_span` — TGP-04), not the derivation basis (what the fold used — TGP-02).
Four other terms remain separate, and Q-β does not touch them.

**QBR-06 (Observed — why it is now a classification question only).** Post-Q-α both P-G0 and P-G1 carry the same
artifact, so Q-β's answer "affects terminology and classification, not whether P-G0 and P-G1 retain different
information" (PGA-19…PGA-21; PGM-31). It is not a behavior-selection question.

**QBR-07 (Observed — nothing is being decided here).** This record sets out what a decision would have to settle. It
does not settle it, and PGM-06's boundary holds: a reading cannot be converted into a decision by analysis.

---

## 2. The two answers, stated precisely

**QBR-08 (Calculated — Answer A: it is also a ground).** The rationale is a **ground** — specifically, the ground of
the transition whose cause it reports. On this answer the artifact carries two classifications at once: a cause (what
produced the event) **and** a ground (what basis was claimed for the change). This is the continuation of ODC-05/ODC-06's
double reading, extended from *report of rationale* to *ground*.

**QBR-09 (Calculated — Answer B: it remains a cause plus a report).** The rationale is **not** a ground. The artifact
is a cause that also carries the actor's stated rationale (ODC-06), and that is all. The corpus has no ground on the
report-only subset, and "the baseline retains no ground" is true as stated.

**QBR-10 (Calculated — the difference between the answers is not about behaviour).** Both attach to the **same
artifact**; neither requires a second field, a marker, a carrier or a new record; neither adds or removes an
obligation. They differ in what the artifact **is called** and in which existing constraint is said to attach to it
(§5).

**QBR-11 (Observed — a third answer is not on the table).** The owner's question is binary and this record keeps it
binary. Nothing here proposes a middle term, a graded classification, or a per-case rule.

---

## 3. What supports Answer A

**QBR-12 (Source-derived — the corpus's own definition of a ground).** GRP-09: a ground "answers **what basis was
claimed for this change**. It is first-person testimony about the act, authored by the process or party whose act it
explains." TGP-02's table fixes the same question: "what basis was claimed for this change?"

**QBR-13 (Source-derived — what the rationale is).** ODC-03: "what the actor says motivated their action."

**QBR-14 (Calculated — the definitional route yields A).** A statement of what the actor says motivated their action,
authored by that actor, about that act, is first-person testimony about the basis claimed for the change. On the
corpus's own definition it **satisfies** the ground's description. CGE-25 reached the same point from the content side:
"the recorded cause is the actor saying why the act happened — which is the same **content** a ground would carry."

**QBR-15 (Inferred — why the definitional route is strong but not decisive).** It establishes that the **content** is
ground-like; it does not establish that "ground" is a term for content. The corpus's usage sometimes treats "ground" as
a **slot in the retention decision** (the object of DQ-2), and on that usage a thing can be ground-like without being
"a ground" in the decision's vocabulary (QBR-17).

**QBR-16 (Calculated — what A would and would not imply).** A implies that the baseline **already retains a ground**
over the report-only subset. It does **not** imply that grounds must be retained anywhere else, that the report must be
marked, that its sufficiency be assessed, or that any position changes (QBR-27). Those are separate questions with
their own answers (QBR-33).

---

## 4. What supports Answer B

**QBR-17 (Calculated — the positional route).** "Ground" in this investigation has done double duty: a kind of
content, and the object of the retention question DQ-2. On the second usage, calling the rationale a ground would say
something about what is **required** — and nothing is required beyond D-3b's cause. Post-Q-α that distinction has no
behavioral consequence, so the label buys nothing and risks implying a requirement that does not exist.

**QBR-18 (Calculated — the caution behind B).** "Ground" is **not canon vocabulary** and no provision requires one
(TGP-05; GRP-08). Naming an existing artifact "a ground" therefore adds a term to a question the corpus has been careful
to keep undefined; B keeps the vocabulary where the text left it.

**QBR-19 (Inferred — and B matches how the decision was written).** ODC-06 chose "recognised as a report of that
actor's stated rationale" — a report, not a basis; and ODC-07 constrained what it establishes rather than what it is.
Read literally, the decision stopped one step short of calling it a ground.

**QBR-20 (Calculated — what B does and does not imply).** B implies that "no ground is retained" is accurate for the
baseline. It does **not** imply that the rationale is worthless, untrustworthy, or unreadable — ODC-07's boundary and
I17 govern its use either way (QBR-25).

---

## 5. What each answer would change

**QBR-21 (Calculated — the complete list of differences).**

| Consequence | **Answer A** (also a ground) | **Answer B** (cause + report) |
|---|---|---|
| **What the artifact is called** | a cause **and** a ground | a cause carrying a report of rationale |
| **Does "the baseline retains no ground" hold?** | **no** — it retains one over the report-only subset | **yes** |
| **Which constraint is said to attach** | CLV-28 — "referenced by the transition, never read as the status" | ODC-07 — the report establishes no truth, authority or validity |
| **Is P-G0's *name* accurate?** | **no**, over that subset | **yes** |
| **How DQ-2's answer is described** | "no ground requirement beyond the one already present" | "no ground requirement at all" |
| **Is K13/K15/K17 a closer analogy?** | slightly — a recorded, attributed basis whose sufficiency is not assessed | no change; analogy only, transferring nothing (GRP-46/47) |
| **Retained content, permissiveness, axis H, any requirement** | **unchanged** | **unchanged** |

**QBR-22 (Calculated — the constraint question is nearly a wash).** CLV-28 is an **inference**, not canon, and D-3b
carries no "never read as warrant" clause (CGE-39; ODC-21). ODC-07 supplies an owner-level boundary covering the same
ground. So A attaches an inferred constraint; B attaches an owner-stated one. Neither is stronger in force, and both
say the same thing in practice: what is recorded is a claim, never a standing.

**QBR-23 (Calculated — so the decision is about a name, and about which sentence describes DQ-2).** That is not
trivial — a reader of a lifecycle record should be able to tell what the recorded cause **is** (ODC-22's seam) — but it
is also not a retention decision, and its consequences are confined to QBR-21's table.

**QBR-24 (Observed — neither answer touches the prohibition question).** P-G0 states no prohibition (PGM-29). Answer A
would make "no ground" inaccurate over a subset; it would **not** create a prohibition on retaining one, and B would not
create a permission to ignore one. The prohibition remains a separate choice requiring its own decision (PGM-30).

**QBR-25 (Observed — neither answer changes what may be done with the artifact).** I17 governs either way: authority
comes from the conditions of establishment, not from the record (`:305`). A report classified as a ground is still a
report; a report not so classified is still subject to the same prohibition on substitution (TGP-19; B-3).

**QBR-26 (Unestablished — the size of what A would cover).** If A is chosen, how much of the stream carries a ground
depends on the report-only subset's extent, which is undrawn (CGE-29/CGE-30). Post-Q-α that size no longer separates
positions (PGA-23); it only bounds the description.

---

## 6. What neither answer can change

**QBR-27 (Observed — the invariant set).**

| Invariant | Because |
|---|---|
| **Retained content under every position** | M-1 binds statement-less events regardless (QAD-05; PGA-13) |
| **Number of requirements** | Q-β classifies an existing obligation; it adds none (ODC-11; QBR-33) |
| **Axis H — the validity gap** | no answer supplies a condition of admissibility (GRP-31) |
| **ND-1** | who may cause a change is untouched (QAD-12) |
| **The validity predicate; the admissible-transition set** | absent, and independent (CEV-24; LCM-24) |
| **QAD-19 — which transitions are statement-less** | a different boundary, still undrawn |
| **P-G2's status** | the only position that adds an obligation (GRP-27) |
| **The prohibition option** | still a separate decision (QBR-24; PGM-30) |
| **ODC's four-term distinction** | cause, rationale, authority and validity stay distinct (ODC-03) |

**QBR-28 (Calculated — one terminological precision, offered so the eventual write-up is exact).** P-G2 would require a
**stated basis recorded on every transition** — a record of **claims**, not of conditions (GRP-35). Describing it as
"retain grounds" is loose about kind but accurate about coverage; the distinction matters because a recorded claim is
what a ground is, and what it is not is a warrant (TGP-03).

**QBR-29 (Observed — nothing is canonical).** This record, and every record in this sequence, is **research evidence**:
no canonical capture, no ledger entry, and no authority to be cited as policy. That status is unchanged by anything
here and would change only by separate authorisation.

---

## 7. What a decision would have to settle

**QBR-30 (Calculated — stated as a checklist of terms, with none filled in).** A Q-β decision would need to fix:

| # | Term | Why it must be explicit |
|---|---|---|
| **T-1** | the answer, A or B | it is not derivable; the two routes in §3/§4 are both available (QBR-15/QBR-17) |
| **T-2** | the scope — the **report-only** case, and that observable causes are excluded | otherwise the answer is read as "all causes", contradicting CGE-26 (QBR-02/QBR-03) |
| **T-3** | that it is a **classification**, creating no requirement, field, marker or carrier | the standing boundary in every prior record (ODC-09; PGM-35) |
| **T-4** | that it creates **no authority** and settles neither ND-1 nor the validity question | the slide QBR-25 forecloses |
| **T-5** | that it does **not** create a prohibition, and does not answer whether one should exist | QBR-24/PGM-30 |

**QBR-31 (Observed — T-2 is the one that is easy to lose).** The question is asked about reported rationales; an answer
recorded without its scope would be read as a claim about causes generally, which the corpus already contradicts
(CGE-26).

**QBR-32 (Observed — no recommendation).** This record does not indicate which answer is preferable, does not predict
which the owner will choose, and does not treat either as more consistent with the prior decisions. Both are consistent
with them (QBR-15/QBR-19).

---

## 8. Traceability

**QBR-33 (Observed).**

| What | Where |
|---|---|
| The question, and the four-term distinction it sits inside | ODC-01, ODC-03, ODC-05…ODC-07 |
| Q-β named as the residual | ODC-18, ODC-19, ODC-22 |
| Report-only vs observable causes | CGE-24, CGE-25, CGE-26, CGE-27 |
| Report-only ≠ statement-less; both undrawn | CGE-29, CGE-30 |
| The ground's definition | TGP-02, TGP-03; GRP-09 |
| "Ground" is not canon | TGP-05; GRP-08 |
| A ground is self-reported; establishes nothing | TGP-20; GRP-09 |
| CLV-28's constraint; D-3b's missing clause | CLV-28; CGE-39; ODC-21 |
| The promotion precedent transfers nothing | GRP-46, GRP-47; K13/K15/K17 as analogy |
| Q-β no longer discriminates positions | PGA-19…PGA-23, PGM-31 |
| P-G0 states no prohibition | PGM-29, PGM-30 |
| Positions' invariant content | PGA-13, PGA-14, PGA-27 |
| Axis H | GRP-31 |
| D-3b's reach, and per-event scope | QAD-04…QAD-08, QAD-19 |

---

## 9. Limits and record status

**QBR-34 (Observed — method).** **Newly read for this record:** TGP-02/03/04 verbatim, CGE-24…CGE-35 verbatim, GRP-09,
and the ODC and PGM records' classification entries. **Not inspected:** any implementation's cause or basis handling,
r32's bytes, the frozen preregistration, v9 beyond §7. No test, probe, benchmark or model call was run (QBR-35).

**QBR-35 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The question is a classification of content; it needs no execution.

**QBR-36 (Observed — the labels' limits).** QBR-02, QBR-03, QBR-08…QBR-10, QBR-14, QBR-21…QBR-23 and QBR-30 are
**Calculated** — read from ODC/CGE/TGP/GRP/PGA/PGM, not from canonical text. QBR-26 and QBR-04 are **Unestablished** on
purpose. Nothing here is offered as a provision of r33 or v9; "ground" has no canonical standing at all.

**QBR-37 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Q-β answered, or either answer ranked | **no** |
| Obligations, prohibitions or carriers created | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence**, no canonical capture, no ledger entry |

**End of record.** Q-β is **open** and is a **classification** question: whether the actor-reported rationale is a
ground or a cause carrying a report (QBR-08/QBR-09). Answer A is supported by the corpus's own definition of a ground
(GRP-09) applied to testimony (CGE-25); Answer B by the term's use as a retention slot and by the decision's own wording
(QBR-17/QBR-19). **Neither changes retained content, permissiveness, any requirement, ND-1, the validity question or
QAD-19** (QBR-27). The consequences are confined to the name, the description of DQ-2, and which constraint is said to
attach (QBR-21). **Nothing is answered, ranked or created; this remains research evidence, not policy.**
