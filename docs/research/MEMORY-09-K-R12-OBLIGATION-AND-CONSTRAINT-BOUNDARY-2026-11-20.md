# MEMORY-09 — K-R12 Obligation and Constraint Boundary

## What class/scope changes require, what A-6 and NSD-09 actually bar, and which constraint is binding

**Status:** obligation and decision-boundary analysis. **Research evidence only — no canonical edit, no schema, no
code, no test, no probe, no ledger entry, no implementation, no owner decision request.** **No carrier is selected;
A-6 and NSD-09 are not revised; the 118-row map is unchanged.** C-1 and C-2 are not reopened. Q1 is not answered, and
this record does not assume that the carrier problem determines N-a's scope.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | what exact obligation does K-R12 create for class/scope changes, and which constraint blocks its representation? |
| Finding | K-R12 creates a **conditional universal**: every actual class/scope change must be recorded in the lifecycle stream and projected from it. The current constraints supply **no admissible carrier**, so the class/scope part is **acknowledged but unrepresented** |
| Binding constraint | **carrier-family-dependent**: A-6 for schema-bearing carriers (P-1, P-2, P-6); K-R14 / "of that stream" for the external-join carrier (P-4); NSD-09 withholds authorisation for any of them. **No single constraint is the blocker across the whole space** |
| Does not do | answer C-3 or Q1; choose a carrier; revise A-6/NSD-09; add a requirement; change the map |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Task and carried state

**KRO-01 (Observed — the owner's request).** Trace the **exact obligation** K-R12 creates for class/scope changes,
then compare it with **A-6** and **NSD-09**. The question is not which carrier to choose; it is whether the current
constraints leave an **acknowledged requirement that cannot presently be represented**, and **which constraint**
would need owner reconsideration to make representation possible. Keep the map unchanged.

**KRO-02 (Observed — carried forward).** Q2 = **Yes** (Q2A-33): post-formation scope and class changes are within the
model's semantic domain; occurrence, authorisation, warrant and representation remain separate (owner, 2026-11-20).
Q5 is live, not conditional. Q1 (which claim C-3 is about) remains open. C3D-11/C3D-18/19: no positive carrier is
admissible under the current constraints; the null carrier is buildable but incomplete.

**KRO-03 (Observed — the boundary).** No carrier selection, no A-6 or NSD-09 revision, no new kind/field/key/schema,
no requirement and no map change. This record does not draft a request and does not decide N-a's scope.

---

## 2. K-R12, parsed

**KRO-04 (Source-derived — the sentence, r33 `:111`, R1b `[X]`).**

> "Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus
> `scope` and `knowledge_class` **changes** — and a memory's current status, class and scope are a **projection of
> that stream**, not mutable authoritative fields."

**KRO-05 (Calculated — three clauses, not one).**

| Clause | Content | What it does |
|---|---|---|
| **C-locus** | the lifecycle is its **own** authoritative stream | locates authority in the stream, separate from acquisition (K-R2/K-R14) |
| **C-domain** | four status kinds **plus** scope and class changes | declares the stream's **subject**, including class/scope changes |
| **C-projection** | current status, class and scope are a **projection** of the stream, not mutable fields | requires derivation and forbids direct authoritative field mutation |

**KRO-06 (Calculated — the obligation, stated exactly).** For a memory `m`, a dimension `d ∈ {scope, class}`, and
times `t0 < t1`, let `v_d(m,t)` be the value at `t`. K-R12 entails:

> **O-KR12(d):** for every `m`, `d`, `t0 < t1`, **if** `v_d(m,t0) ≠ v_d(m,t1)`, **then** there is a record
> `r ∈ L(m)` whose content determines `v_d(m,t1)`, and the projection of `L(m)` at `t1` yields `v_d(m,t1)`.

**KRO-07 (Calculated — the logical form matters).** O-KR12 is a **conditional universal**: it quantifies over
**actual** changes and says where they must be recorded and how the current value is obtained. It does **not** assert
that any change occurs, and it does not assert that every possible change must be enabled. It is satisfied
vacuously if no change occurs.

**KRO-08 (Calculated — the capacity question is separate).** C-domain says class/scope changes are **within the
stream's subject**. That supports a stronger reading:

| Reading | Claim | Consequence for the current constraint set |
|---|---|---|
| **strong (capacity)** | the stream must be **able** to carry class/scope changes, because they are declared contents | the defect is **present now**: the stream cannot carry one of its declared subjects |
| **weak (conditional)** | the stream is where such changes **would** be recorded; no change, no record needed | the defect is **trigger-conditional**: compliance holds until the first change |

**KRO-09 (Unestablished — which reading K-R12 requires).** The text supports the strong reading more naturally — a
subject declaration is a claim about what the stream is **for** — but the weak reading cannot be excluded. The
difference is whether the gap is **present-now** or **trigger-conditional**, not whether a gap exists. Both readings
leave the class/scope part unrepresented once a change occurs.

**KRO-10 (Calculated — O-KR12 is not a new requirement).** It is K-R12's existing content, traced from its text. It
adds no row, no predicate and no canonical change. It is **acknowledged in canon**, whether or not the 118-row map
captures it (the map is v9-derived and has no class/scope-change row; KRO-45).

**KRO-11 (Calculated — what O-KR12 does not require).** It does not require that changes occur, that any particular
cause be permitted, that a sufficiency standard exist, that an actor be authorised, or that a particular carrier be
chosen. Those are Q2's separate questions and Q5's.

---

## 3. What the current vocabulary supplies

**KRO-12 (Observed — status has a vocabulary, class/scope do not).** K-R12 names four status kinds
(`established | superseded | corrected | retracted`). **No kind names a class or scope change** (C3Q-07). The
class/scope part is carried by the phrase "plus … changes," not by a named entry type.

**KRO-13 (Observed — no lifecycle schema exists for any entry).** A2-11/A2-12 and ECI-11: K-R12 is authoritative by
declaration and unspecified by schema. The four status names are **vocabulary without a schema**; class/scope are
**neither vocabulary nor schema**.

**KRO-14 (Calculated — the asymmetry is precise).**

| Part of C-domain | Named entry type? | Schema? | Gap |
|---|---|---|---|
| status transitions | **yes** — four kinds | **no** | vocabulary exists; representation does not |
| scope changes | **no** | **no** | neither exists |
| class changes | **no** | **no** | neither exists |

**KRO-15 (Calculated — why the status gap does not excuse the class/scope gap).** A schema for the four status kinds
would still not carry a class/scope change: the kinds are status-shaped, and K-R12 lists class/scope **after** them,
as a separate subject. A status schema plus a class/scope **attribute** is P-6; a status schema plus a separate
**kind** is P-1/P-2. Neither is the status schema alone.

---

## 4. A-6 and NSD-09, exactly

**KRO-16 (Source-derived — A-6, NAE-31).**

> "**A-6** — that the decision **creates no event kind, field, key or schema**, and does not settle authority or
> validity — the standing boundaries (NSD-09; GRP-31)."

**KRO-17 (Source-derived — what A-6 is).** A-6 is a **term of the N-a fold decision** (NAE-31), not a standalone
canonical rule. Its purpose is to keep a fold decision from becoming a schema decision (NAE-32). It binds **the
decision**, not the architecture forever.

**KRO-18 (Source-derived — NSD-09, owner scope determination C-6).**

> "**NSD-09** — **No substantive authorisation.** The determination identifies questions, not their answers: no fold
> rule, ordering key, vocabulary mapping, event, actor, predicate or schema."

**KRO-19 (Source-derived — what NSD-09 is).** NSD-09 is a term of the **ND-1 scope determination**: that
determination identifies questions and authorises none of their answers. It is not itself a schema rule; it is an
**authorisation boundary** on what that determination may settle.

**KRO-20 (Calculated — the two are not the same kind of constraint).**

| Constraint | Binds | Says | Does not say |
|---|---|---|---|
| **A-6** | the **N-a decision** | this decision creates no kind, field, key or schema | that no future decision may |
| **NSD-09** | the **ND-1 determination** | this determination authorises no substantive answer | that the requirement is invalid or permanently unanswerable |
| **standing read-only boundary** | **this record** | no schema, carrier, field, event type or mechanism | that the owner may not authorise one elsewhere |

**KRO-21 (Calculated — neither denies O-KR12).** A-6 bars **this decision** from supplying a carrier; NSD-09 bars
**this determination** from authorising one. Neither says class/scope changes are impossible, nor that K-R12's
obligation is withdrawn. **The requirement remains; the means are withheld.**

---

## 5. Comparison: obligation against constraints

**KRO-22 (Calculated — the direct conflict is carrier-family-specific).**

| Carrier family | What it needs | Which constraint binds | What reconsideration would mean |
|---|---|---|---|
| **P-1** own transition kind | a new lifecycle kind | **A-6** ("no event kind") | the N-a decision, or a separate authorised decision, would have to permit a kind |
| **P-2** generic change entry | a new entry kind | **A-6** ("no event kind") | same |
| **P-6** piggyback on an existing entry | a new attribute/field | **A-6** ("no field") | a field would have to be authorised; the narrowest schema-bearing change |
| **P-4** external operation joined at projection | another operation as input | **K-R14** ("rebuilt independently") and "of that stream" | A-6 is irrelevant; the architecture's independence rule would have to change |
| **P-5** stored projection | an authoritative input and rebuild rule | **none** — it is storage | no constraint change makes it a carrier; it needs the input P-1…P-4/P-6 would supply |
| **P-3** null carrier | nothing new | **none** | compliant, but O-KR12 is unmet whenever a change occurs |

**KRO-23 (Calculated — no single constraint is the blocker).** A-6 blocks the schema-bearing families; K-R14 blocks
the external-join family; NSD-09 withholds authorisation for any of them; P-5 is not a carrier and P-3 accepts the
gap. **"Which constraint" has no unique answer until the carrier family is fixed** — and fixing it is Q5, which is
not this record's task.

**KRO-24 (Calculated — A-6 is the operative bar only inside the N-a frame).** If the owner opened a **separate,
authorised representation decision**, A-6 would not apply to it (it binds the N-a decision). The obstacle would then
be the absence of authorisation, not A-6's text. This is why "reconsider A-6" and "authorise a separate decision" are
different remedies, and why the record does not recommend either.

**KRO-25 (Calculated — P-4 is not an A-6 workaround).** P-4 needs no new kind or field, so A-6 does not block it —
but K-R14 and "of that stream" do. Reconsidering A-6 would not make P-4 admissible. The constraint that binds is the
**independent-rebuild rule**, not the schema bar (C3D-13).

**KRO-26 (Calculated — P-6 is the narrowest schema-bearing change).** It needs an attribute, not a kind. That is a
statement about the **size** of the constraint change, not a recommendation: P-6 is not selected, and its textual
adequacy ("piggyback on a status entry") remains a Q5 question.

---

## 6. Is there an acknowledged requirement that cannot presently be represented?

**KRO-27 (Calculated — yes, for the class/scope part).** O-KR12(d) is acknowledged in K-R12. The current constraint
set supplies no admissible carrier for a class/scope change (C3D-11). Therefore:

> **For `d ∈ {scope, class}`, if a post-formation change occurs, K-R12's recording-and-projection obligation cannot
> be satisfied with the current vocabulary** — because there is no record type that can carry the change into the
> lifecycle stream, and no admissible external input under K-R14.

**KRO-28 (Calculated — the two readings give two forms of the same finding).**

| Reading | Present state | When it becomes a violation |
|---|---|---|
| **strong (capacity)** | the stream cannot carry a declared subject — **unrepresented now** | already |
| **weak (conditional)** | compliant while no change occurs; **unsatisfiable on the first change** | at the first class/scope change |

**KRO-29 (Calculated — the status part is different).** The status part has named vocabulary but no schema; its gap
is a **schema gap**. The class/scope part has neither vocabulary nor schema; its gap is a **vocabulary-and-schema
gap**. Both are unrepresented at the schema level, but only the class/scope part requires **new vocabulary** (a kind
or a field) in addition to a schema.

**KRO-30 (Calculated — this is not a claim that a change has occurred).** The finding is about **capacity and
obligation**, not occurrence. It does not assert that any memory's class or scope has changed; it asserts that if one
does, the current representation cannot record it as K-R12 requires.

**KRO-31 (Observed — the map does not capture it).** The 118-row map has no class/scope-change requirement; its
scope-related row is L7 (authority grants), not the memory `scope` field. The requirement is acknowledged **in
canon**, not in the map. That is a capture gap, not a denial (DQ-3; KRO-10).

---

## 7. Which constraint would need owner reconsideration?

**KRO-32 (Calculated — the answer is conditional, and that is the precise result).**

| If the owner were to pursue… | The constraint to reconsider | Why that one |
|---|---|---|
| a new lifecycle kind (P-1/P-2) | **A-6** inside N-a, or a **separate authorised decision** outside A-6 | A-6's "no event kind" is the operative bar |
| a field on an existing entry (P-6) | **A-6** ("no field") | narrowest schema-bearing change; still barred |
| an external operation (P-4) | **K-R14** / "of that stream" | A-6 does not block it; independence does |
| a stored projection only (P-5) | **none** — it is not a carrier | it needs an authoritative input, which is the carrier question |
| no change (P-3) | **none** | compliant, but O-KR12 stays unmet on a change |

**KRO-33 (Calculated — NSD-09 is common to every positive route).** Any positive carrier requires an **authorised
decision**. NSD-09 currently withholds that authorisation from the ND-1 determination. So the common constraint is
not A-6 but the **authorisation boundary**; A-6 is the operative schema bar only for the schema-bearing families.

**KRO-34 (Calculated — "reconsider A-6" is not the same as "revise the architecture").** A-6 is a decision term. The
owner could (a) reconsider it for N-a, (b) authorise a separate representation decision where it does not apply, or
(c) leave it and accept the gap. This record identifies the options; it does not choose among them and does not
recommend (a) or (b).

**KRO-35 (Calculated — no constraint change is authorised here).** Naming which constraint binds is a finding about
the decision boundary, not a request to move it. The read-only boundary stands; the map stays at 118.

---

## 8. The Q1 qualification

**KRO-36 (Observed — the owner's condition).** Q1 still determines whether the next decision concerns **N-a's fold
scope**, the **lifecycle stream's broader representation**, or **K-R12's projection requirement**. These are related
but not interchangeable. Investigating the positive-carrier problem under K-R12 does not determine N-a's scope.

**KRO-37 (Calculated — O-KR12 attaches to K-R12, not to N-a).** The obligation is a property of the **stream and its
projection**. N-a is one possible **discharge site** — the fold that produces current state — but it is not the
obligation's source. A carrier could be located in the broader lifecycle representation, or the projection
requirement could be stated without a fold at all.

**KRO-38 (Calculated — A-6's bar is a property of the N-a decision, not of the obligation).** The fact that A-6
prevents **N-a** from creating a carrier does not show that the obligation is unrepresentable **in principle**; it
shows that N-a, as currently bounded, cannot discharge it. Whether that makes N-a status-only, or moves the carrier
elsewhere, is Q1's question.

**KRO-39 (Calculated — the correct conclusion is conditional).** "The current constraints leave an acknowledged
requirement that cannot presently be represented" is true **for the class/scope part** and **under the current
decision frame**. It does **not** entail "N-a must be status-only," "the architecture is inconsistent," or "a new
kind is required." Those are Q1/Q5 conclusions this record does not draw.

---

## 9. What this record does not do

**KRO-40 (Observed).**

| Untouched | State |
|---|---|
| **Q1 / C-3** | open; no scope determination |
| **Q5** | live; no carrier selected |
| **Q4** | premature |
| **A-6 / NSD-09** | not revised, not relaxed |
| **carrier** | none selected; no positive carrier admissible today |
| **requirement / map** | unchanged; **118** rows |
| **N-a / N-b / N-c / entry architecture** | untouched |
| **C-1 / C-2** | determined / non-decision; not reopened |
| **§11 adoption; §8.1 collision; O-8a** | untouched |
| **owner request** | not drafted |

**KRO-41 (Observed — no new rule, predicate, kind, field or schema).** O-KR12 is a **paraphrase of K-R12's existing
text**, not a new requirement. It creates no canonical record.

---

## 10. Traceability

**KRO-42 (Observed).**

| What | Where |
|---|---|
| K-R12: locus, domain, projection, non-mutation | `:111`; normative extract `:106` |
| K-R14: independent rebuild; no stored join | `:113` |
| K-R13: cross-stream read-time join | `:112` |
| A-6: no kind/field/key/schema; does not settle authority or validity | NAE-31; NAE-32 |
| NSD-09: no substantive authorisation | owner scope determination C-6 |
| No lifecycle schema; declaration without schema | A2-11/A2-12; ECI-11 |
| No kind names a class/scope change | C3Q-07 |
| Null vs positive carrier; P-5 is storage | C3D-05…C3D-11 |
| P-4's independence strain | C3D-13 |
| P-6 is the narrowest positive carrier | C3D-14 |
| Q2 = Yes; domain vs representation | Q2A-33…Q2A-35 |
| Q1's three framings; no N-a inference | owner, 2026-11-20; KRO-36…KRO-39 |
| Map has no class/scope-change row | 118-row map; L7 is authority-scoped |

---

## 11. Limits and record status

**KRO-43 (Observed — method).** **Newly read:** K-R12/K-R13/K-R14 (`:111`–`:113`); NAE-31/32; NSD-09; C3D-05…C3D-19;
C3Q-06…C3Q-09; Q2A-33…Q2A-42; the map's scope rows. **Not inspected:** any implementation's lifecycle handling, r32's
bytes (absent), the frozen preregistration, v9 beyond §7.

**KRO-44 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. This is a textual and decision-boundary analysis.

**KRO-45 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier selected; A-6/NSD-09 revised | **0** |
| Q1/C-3 answered; request drafted | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** K-R12 creates **O-KR12(d)**: every actual class/scope change must be recorded in the lifecycle
stream and projected from it (KRO-06). The current vocabulary has **no entry type** for such a change — status has
four names but no schema; class/scope have neither (KRO-12…KRO-15). **A-6** bars the **N-a decision** from creating a
kind, field, key or schema; **NSD-09** withholds substantive authorisation from the ND-1 determination (KRO-16…KRO-21).
The binding constraint is **carrier-family-dependent**: A-6 for P-1/P-2/P-6, K-R14 for P-4, none for P-3/P-5
(KRO-22…KRO-26). **The acknowledged requirement is therefore unrepresented for the class/scope part** — presently
under the capacity reading, on the first change under the conditional reading (KRO-27…KRO-30). **No single constraint
is the blocker, and no constraint change is authorised** (KRO-32…KRO-35). O-KR12 attaches to **K-R12**, not to N-a;
**Q1 still decides whether N-a's fold, the broader stream representation, or the projection requirement is the next
site** (KRO-36…KRO-39). **Nothing is selected, requested, designed or authorised.**
