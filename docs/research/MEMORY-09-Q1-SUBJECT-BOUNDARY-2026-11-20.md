# MEMORY-09 — Q1 Subject Boundary

## Which decision frame owns the class/scope carrier problem

**Status:** subject-boundary analysis. **Research evidence only — no canonical edit, no schema, no code, no test, no
probe, no ledger entry, no implementation, no owner decision request.** **No carrier is selected; no architecture is
proposed; no requirement is captured; the 118-row map is unchanged.** C-1 and C-2 are not reopened. Q2 = Yes is
carried forward; Q5 is owned but not answered; Q4 is narrowed but not decided.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | which of the three Q1 subjects — N-a's fold scope, the lifecycle representation outside N-a, or K-R12's projection requirement — owns the unresolved carrier problem? |
| Finding | **the broader lifecycle representation (S-2)** owns it. **S-3** (the projection requirement) is the acceptance test, not a decision frame. **S-1** (N-a's fold scope) is a downstream consumer that A-6 bars from creating the carrier |
| Alternative | expanding N-a to include the carrier would require reopening A-6; the evidence does not support that as the first move |
| Does not do | answer C-3 as a single question; select a carrier; propose a schema; capture a requirement; change the map |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Task and carried state

**Q1S-01 (Observed — the owner's request).** Resolve Q1's subject boundary: map the three subjects against their
existing authority, constraints and unresolved questions; decide which frame owns the next question. No new
architecture, no carrier selection, no requirement capture.

**Q1S-02 (Observed — carried forward).** Q2 = **Yes** (Q2A-33): class/scope changes are within the semantic domain;
occurrence, authorisation, warrant and representation are separate. K-R12 creates **O-KR12(d)** — every actual change
must be recorded in the lifecycle stream and projected from it (KRO-06). The current vocabulary has **no entry type**
for a class/scope change; the null carrier is buildable but incomplete (C3D-11; KRO-27).

**Q1S-03 (Observed — the boundary).** No carrier, no schema, no kind, no field, no requirement, no map change, no
architecture. This record maps **who owns the question**, not what the answer is.

**Q1S-04 (Observed — Q1 as already decomposed).** C3Q-14 stated Q1 as: is C-3 a question about **N-a** (does the fold
cover class/scope?), about the **stream** (does the lifecycle need a class/scope representation?), or about the
**projection requirement** (does K-R12 require a carrier even without a fold?). This record keeps those three and
does not add a fourth.

---

## 2. The three subjects, defined

**Q1S-05 (Source-derived — S-1, N-a's fold scope).** N-a is "what a later lifecycle event does to the current state —
the transition function over K-R12's kinds" (NDR). Its terms are A-1 (shape), A-2 (domain), A-3 (entry), A-4
(supersession links), A-5 (order sensitivity), A-6 (no kind/field/key/schema; no authority or validity) (NAE-31).

**Q1S-06 (Source-derived — S-2, the lifecycle representation outside N-a).** The stream as an object: its entry
vocabulary, schema, entry architecture and projection, considered apart from any particular fold. K-R12 declares the
stream's locus and domain (`:111`); no schema exists for any lifecycle entry (A2-11/A2-12; ECI-11); K-R14 requires
independent rebuild (`:113`).

**Q1S-07 (Source-derived — S-3, K-R12's projection requirement).** The canonical requirement that current status,
class and scope be a **projection** of the stream, not mutable authoritative fields (`:111`). Its formal content is
O-KR12(d) (KRO-06). It constrains any representation; it does not create one.

**Q1S-08 (Calculated — the three are not nested).** S-1 is a **consumer** of the stream; S-2 is the **object**; S-3 is
a **criterion** on the object. A decision about one does not automatically decide the others (the owner's condition,
KRO-36).

---

## 3. Authority map

**Q1S-09 (Calculated — who can decide what).**

| Subject | Existing authority | What that authority can do | What it cannot do |
|---|---|---|---|
| **S-1** N-a fold | the owner's N-a decision, bounded by A-1…A-6 | fix a fold over **existing** kinds; choose shape, entry, links, order | create a kind, field, key or schema; settle authority or validity (A-6) |
| **S-2** stream representation | **no standing decision**; would require explicit owner authorisation | create/decline a class/scope representation; choose its form; define its projection | decide N-a's fold, N-b's vocabulary or the authority/validity questions unless separately scoped |
| **S-3** projection requirement | **K-R12 itself**, canonical | state the test any representation must pass | create a carrier; decide a schema; authorise anything |

**Q1S-10 (Source-derived — A-6 binds S-1, not S-2 by default).** A-6 is a term of the **N-a decision**: "the decision
creates no event kind, field, key or schema, and does not settle authority or validity" (NAE-31). It is not a
freestanding prohibition on the architecture (NAE-32; KRO-17). A separate representation decision is not bound by it
unless the owner imports it.

**Q1S-11 (Source-derived — NSD-09 binds the ND-1 determination, not S-2).** NSD-09 says the **ND-1 scope
determination** identifies questions and authorises no substantive answer. It is not a schema rule and not a
permanent bar on a separately authorised decision (KRO-19). It is the common **authorisation boundary**: identifying
a possible route does not authorise taking it.

**Q1S-12 (Observed — the read-only boundary binds this record).** No schema, carrier, field, event type or mechanism
is authorised **here**. That is a constraint on the research, not a statement that the owner may never authorise
S-2.

**Q1S-13 (Calculated — the authority chain, stated precisely).** S-1 is bounded by A-6; S-2 is bounded by the need for
explicit authorisation; S-3 is bounded by nothing new — it is already canonical. **The three have different owners and
different gates.** Treating them as interchangeable is the error this map prevents.

**Q1S-14 (Observed — no new authority is created).** Naming S-2 as the owner of the carrier problem does not open it,
authorise it, or make it a decision. It identifies which frame **would** own it if the owner opened one.

---

## 4. Constraint map

**Q1S-15 (Calculated — what each subject must respect).**

| Subject | Binding constraints | Not binding by default |
|---|---|---|
| **S-1** | A-6; K-R12's projection; K-R14 independent rebuild; the 118-row map; the read-only boundary | — |
| **S-2** | K-R12 (locus, domain, projection); K-R14 (independent rebuild; no stored join); K-R2 (rebuild/integrity); `:366` check; §11 authority (C-1 R-E), distinct | **A-6** (N-a term); **NSD-09** (ND-1 term) — unless imported |
| **S-3** | K-R12; K-R14; `:366`; `:367` projection isolation | any schema decision; it has none of its own |

**Q1S-16 (Calculated — A-6's reach is the crux).** A-6 bars **S-1** from creating the carrier. It does **not** bar
S-2, because S-2 is not the N-a decision. This is why "the class/scope part is unrepresented" does not entail "the
architecture forbids representing it" (KRO-24, KRO-38).

**Q1S-17 (Calculated — K-R14's reach is different).** K-R14 binds **any** representation: each stream is rebuilt
independently, and the cross-stream join has no stored form (`:113`). It is why P-4 is not an A-6 workaround (KRO-25).
S-2 must satisfy it; S-1 cannot waive it.

**Q1S-18 (Calculated — S-3 is already satisfied in form, not in content).** The projection requirement is canonical
and applies now. Its **content** for class/scope is settled by Q2 = Yes; its **implementation** has nothing to derive
because no carrier exists (KRO-27). S-3 is therefore a **test with an unmet input**, not an open question about what
the test requires.

**Q1S-19 (Observed — the map is unchanged).** No subject adds a requirement. S-2 would decide **how** to satisfy
K-R12's existing obligation, not create a new one (KRO-10). The map stays at 118; DQ-3's capture gap is acknowledged.

**Q1S-20 (Calculated — the read-only boundary blocks all three equally).** No subject is opened by this record. The
difference is **which gate** a future decision would have to pass: A-6 for S-1, authorisation for S-2, nothing new for
S-3.

---

## 5. Unresolved-question map

**Q1S-21 (Calculated — what remains open under each subject).**

| Subject | Unresolved questions | Relation to other subjects |
|---|---|---|
| **S-1** | fold shape (F-1…F-6); A-2 domain (status only vs status+class+scope); A-3 entry; A-4 supersession links; A-5 order sensitivity; N-c conditional on A-5 | N-b (status vocabulary) sits here; N-c is conditional on A-5 |
| **S-2** | Q5 carrier; entry architecture; schema; kind/field; projection rule; check; ordering for class/scope entries; warrant/authority for changes; relation to §11 and O-8a | S-3 is its acceptance test; S-1 is a downstream consumer |
| **S-3** | no projection rule for class/scope because no carrier; no check for class/scope; depends on S-2 | settled in content by Q2 = Yes; unmet in implementation |

**Q1S-22 (Calculated — S-1's A-2 is not free).** A-2 asks whether the fold's domain is status only or
status+class+scope. The second option presupposes a class/scope entry; none exists, and A-6 bars S-1 from creating
one. So **A-2 is decidable only as "status only" within S-1's frame** — and that is a statement about S-1's
vocabulary, not about the semantic domain (Q2 = Yes).

**Q1S-23 (Calculated — S-2's questions are the carrier questions).** Q5, entry architecture, schema and projection
rule are all S-2 questions. They are not N-a questions, because A-6 bars N-a from answering them. This is the
mapping's central separation.

**Q1S-24 (Calculated — S-3 has no independent decision).** Once Q2 = Yes, S-3 says "a carrier is required." It does
not say which. Its only remaining role is to **test** an S-2 proposal. It cannot own the next question because it has
no authority to create anything (Q1S-09).

**Q1S-25 (Observed — N-b and N-c stay with S-1).** N-b (status vocabulary: `active`/`established`) is a fold-input
question; N-c (ordering key) is conditional on A-5's order sensitivity (NSD-11/12). Neither is moved by S-2.

**Q1S-26 (Calculated — S-2 may raise a new ordering question, not N-c).** If a class/scope entry is created, it may
need its own ordering rule. That would be a **new** S-2 question, not N-c, which is about the status stream. Naming
this is not opening it.

**Q1S-27 (Observed — the entry architecture is S-2's, not N-a's).** The entry architecture concerns how lifecycle
entries are shaped and related; it is part of the stream's representation. N-a's A-3 asks only where the fold starts,
not how entries are built.

---

## 6. Can each subject own the carrier problem?

**Q1S-28 (Calculated — S-1 cannot).** Three reasons: (i) A-6 bars it from creating a kind, field, key or schema;
(ii) deciding "status only" would convert a **vocabulary limitation** into a **semantic conclusion**, which the owner
has repeatedly refused; (iii) deciding "status+class+scope" would be unimplementable without a carrier. S-1 can own
the **fold over status**; it cannot own the class/scope carrier.

**Q1S-29 (Calculated — S-2 can).** Its scope is the stream's representation; a separate authorised decision can
create a carrier; A-6 does not apply by default; it must satisfy K-R12/K-R14/K-R2. It is the **only** frame whose
scope includes the carrier without importing A-6.

**Q1S-30 (Calculated — S-3 cannot).** It is a requirement, not an authority. It can test a carrier; it cannot create
one. Its content is already settled by Q2 = Yes.

**Q1S-31 (Calculated — the finding).** **S-2 owns the unresolved carrier problem.** S-3 is its acceptance test. S-1 is
a downstream consumer that must not be asked to answer it. This is a **subject-boundary finding**, not a carrier
selection, not a schema proposal, and not an authorisation (Q1S-14).

**Q1S-32 (Calculated — the alternative, stated fairly).** The owner could expand N-a to include the carrier. That
would require reopening A-6 and would make N-a a schema decision. The evidence does not support that as the first
move: it conflates the fold with the representation, and it asks A-6 to do work it was designed to prevent (NAE-32).
This record does not choose between S-2 and an expanded N-a; it records that **as currently constituted**, S-1 cannot
own the carrier.

**Q1S-33 (Calculated — why S-2 is not "new architecture").** S-2 is the **existing** stream described by K-R12; the
question is its representation, not a new design. Mapping it as the owner adds no architecture, no schema and no
requirement (Q1S-19).

---

## 7. What this does to C-3

**Q1S-34 (Calculated — C-3 decomposes rather than resolves).** C-3 ("is N-a status-only?") has three components, and
Q1 assigns them:

| Component of C-3 | Owner | State |
|---|---|---|
| **semantic** — are class/scope within the lifecycle domain? | Q2 / K-R12 | **Yes** (Q2A-33) |
| **representational** — what carries a class/scope change? | **S-2** / Q5 | open; carrier not selected |
| **fold-scope** — does N-a's transition function include class/scope? | **S-1** | status-only **within its frame**; class/scope deferred to S-2 |

**Q1S-35 (Calculated — C-3 is not a single decision request).** It does not have one answer at the N-a level. Its
semantic part is settled by Q2; its representation part belongs to S-2; its fold part is a scoped S-1 decision. Asking
N-a to answer C-3 would force it to answer the representation question, which A-6 forbids.

**Q1S-36 (Calculated — N-a's status-only is a frame boundary, not a semantic conclusion).** S-1 may decide a fold
over the existing status kinds, provided it **explicitly defers** class/scope to S-2. That is legitimate: it says
"this fold does not cover a vocabulary that does not exist," not "class/scope never change." The difference is the
one the owner identified between a present asymmetry and a boundary of meaning (C3Q-09).

**Q1S-37 (Calculated — S-1's deferral is a condition, not an option).** If S-1 decided status-only **without**
deferring class/scope to S-2, it would be read as a semantic denial, which Q2 = Yes contradicts. The deferral is what
keeps the frame decision honest.

**Q1S-38 (Calculated — a later carrier changes S-1's input, not its authority).** If S-2 created a class/scope entry,
S-1 would still need a separate amendment to fold it (A-6 would still bar S-1 unless reconsidered). That is a
dependency, not a contradiction: the representation decision comes first; the fold follows if the owner wants it.

---

## 8. What this does to Q4 and Q5

**Q1S-39 (Calculated — Q5 is owned, not answered).** Q5 (carrier) belongs to S-2. It remains open; no carrier is
selected; the current constraints still admit no positive carrier without a decision (C3D-11).

**Q1S-40 (Calculated — Q4 is narrowed, not decided).** Q4 asks whether class/scope are deferred or excluded. The Q1
finding narrows the live options: "excluded" is inconsistent with Q2 = Yes and K-R12's domain; "deferred to S-2" is
the only disposition consistent with the evidence. The **decision** remains the owner's; this record does not make it.

**Q1S-41 (Calculated — what an S-2 decision could legitimately decide).**

| Could decide | Could not decide |
|---|---|
| whether to create a class/scope carrier in the lifecycle stream | N-a's fold shape, entry or order (S-1) |
| its form, subject to K-R14 (kind, field or another in-stream form) | N-b's status vocabulary (S-1) |
| the projection rule and check for class/scope | authority or validity, unless separately scoped |
| the entry architecture for class/scope changes | its own authorisation (the owner's) |
| whether to decline and accept the gap | a requirement capture or a map change (separate step) |

**Q1S-42 (Observed — S-2 is not authorised by this record).** The read-only boundary stands; the map stays at 118.
The owner would have to authorise S-2 explicitly before any representation work begins.

---

## 9. What this record does not do

**Q1S-43 (Observed).**

| Untouched | State |
|---|---|
| **carrier / schema / kind / field** | none selected or proposed |
| **C-3** | decomposed, not answered as a single question |
| **Q4 / Q5** | Q5 owned by S-2; Q4 narrowed, not decided |
| **N-a / N-b / N-c** | open; N-a's status-only is a frame boundary, not a semantic answer |
| **A-6 / NSD-09** | not revised, not relaxed |
| **map / requirements** | **118**; no capture |
| **architecture** | no new architecture proposed |
| **owner request** | not drafted |

---

## 10. Traceability

**Q1S-44 (Observed).**

| What | Where |
|---|---|
| K-R12: locus, domain, projection | `:111`; extract `:106` |
| K-R14: independent rebuild; no stored join | `:113` |
| K-R13: read-time join | `:112` |
| N-a definition; A-1…A-6 | NDR; NAE-31 |
| A-6 binds the N-a decision | NAE-31/32; KRO-17 |
| NSD-09 binds the ND-1 determination | owner scope C-6; KRO-19 |
| No lifecycle schema | A2-11/A2-12; ECI-11 |
| O-KR12; class/scope gap | KRO-06; KRO-27 |
| Carrier families and binding constraints | C3D-10/11; KRO-22/23 |
| Q2 = Yes | Q2A-33 |
| Q1's three subjects | C3Q-14 |
| N-b/N-c separation | NSD-11/12 |
| Status-only as frame, not semantics | C3Q-09; Q1S-36 |

---

## 11. Limits and record status

**Q1S-45 (Observed).** **Newly read:** C3Q-14; NDR's N-a definition; NAE-31/32; NSD-09/11/12; KRO-06/17/19/22/23;
C3D-10/11; Q2A-33. **Not inspected:** any implementation's lifecycle handling, r32's bytes (absent), the frozen
preregistration. 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound, `ExecutionNotAuthorised` at
`conversation.py:177`.

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier selected; A-6/NSD-09 revised | **0** |
| C-3 answered as a single question; request drafted | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** Q1's subject boundary is **S-2, the lifecycle representation outside N-a** (Q1S-31). S-3, the
projection requirement, is the **acceptance test** — its content is settled by Q2 = Yes, but it has no authority to
create a carrier. S-1, N-a's fold, is a **downstream consumer**: A-6 bars it from creating the carrier, and its
status-only scope is a **frame boundary**, not a semantic conclusion (Q1S-22, Q1S-36). **C-3 decomposes**: its
semantic component is Q2 = Yes, its representational component belongs to S-2/Q5, its fold component is a scoped S-1
decision (Q1S-34/35). **Q5 is owned by S-2 but not answered; Q4 is narrowed, not decided** (Q1S-39/40). **No carrier,
schema, requirement, map change or architecture is proposed, and S-2 is not authorised by this record** (Q1S-42).
