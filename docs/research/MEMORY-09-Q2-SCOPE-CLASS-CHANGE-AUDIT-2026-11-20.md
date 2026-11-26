# MEMORY-09 — Q2 Source Audit

## Can `scope` and `knowledge_class` change after formation?

**Status:** source audit and finding. **Research evidence only — no canonical edit, no schema, no code, no test, no
probe, no ledger entry, no implementation, no owner decision request.** It applies the owner's threshold to the
authoritative text. **C-3 remains unanswered; no carrier is selected; A-6 is not relaxed; the 118-row map is
unchanged.** C-1 and C-2 are not reopened.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | can `scope` and `knowledge_class` change after formation? |
| Finding | **Q2 = Yes** — the ratified text explicitly names scope and class changes, and three rules regulate their **causes**; no text fixes the values |
| What it establishes | **possibility within the model** — not requirement, frequency, sufficiency, mechanism or representation |
| Consequence | a **positive carrier is required**; under the current constraints none is admissible, so the null carrier is buildable but **incomplete** (C3D-11, C3D-18/19) |
| Does not do | answer C-3; choose a carrier; draft a request; revise A-6 or the map |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. The task and the threshold

**Q2A-01 (Observed — the owner's scope).** A source-focused audit of Q2, looking in three places: normative text,
original design intent, and existing behavioural examples. Evidence of **possible change** is kept separate from
evidence of **how change is represented** (Q2A-02).

**Q2A-02 (Observed — the separation).** "Can change" (Q2) is a semantic question. "Has a carrier" (Q3/Q5) is a
representation question. An absent kind is not evidence of immutability, and the word "changes" is not proof that
any particular change is required (C3Q-08/09; owner, 2026-11-20).

**Q2A-03 (Source-derived — the threshold).**

| Finding | Result |
|---|---|
| explicit evidence that changes **can** occur | **Q2 = Yes**; a positive carrier is required |
| explicit evidence that the values are **fixed** after formation | **Q2 = No**; the null carrier may be complete |
| wording ambiguous, or evidence conflicts | **Q2 = Unestablished**; the carrier analysis stays conditional |

**Q2A-04 (Observed — no other outcome).** The audit does not treat silence, a missing kind or a missing example as a
finding of "No." Those are representation facts, not semantic ones.

---

## 2. Normative text

**Q2A-05 (Source-derived — K-R12, r33 `:111`, under R1b `[X]`).**

> "Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus
> `scope` and `knowledge_class` **changes** — and a memory's current status, class and scope are a **projection of
> that stream**, not mutable authoritative fields."

**Q2A-06 (Calculated — what the clause says about change).** It does not merely say the values are projected. It
names **scope and knowledge_class changes as contents of the stream**, in the same sentence that makes the values a
projection of that stream. A "change" to a current value is a transition from a prior value; for a record, that is a
post-formation change. **The clause therefore provides for such changes.**

**Q2A-07 (Source-derived — K-R10, r33 `:224`, under R3 `[X]`).**

> "An experiential event never downgrades to ambient by repetition or by age; decay does not change class (D5)."

**Q2A-08 (Calculated — K-R10 regulates a cause, not the possibility).** "Never downgrades … by repetition or by age"
forbids two **causes**. It does not say class cannot change; it says these things do not change it. The rule is
intelligible only if a class change is a coherent event in the model.

**Q2A-09 (Source-derived — K-R9, r33 `:223`, marked `[D]`).**

> "No record becomes ambient because many characters know it. Popularity is not authorship."

**Q2A-10 (Calculated — K-R9 also presupposes the transition).** "Becomes ambient" is contemplated; what is refused is
**popularity** as the cause. "Popularity is not authorship" leaves authorship as the relevant alternative — without
saying that any re-authoring is permitted. K-R9 is `[D]` (needs review), so it corroborates rather than establishes.

**Q2A-11 (Source-derived — D5, r33 `:72`, owner decisions table).**

> "**Decay** affects surfacing only — never truth, ownership, **scope**, eligibility or history."

**Q2A-12 (Calculated — D5 covers the other dimension, again by cause).** D5 forbids **decay** from affecting scope; it
does not say scope is fixed. Read with K-R10, the corpus regulates the causes of change in **both** dimensions —
class by repetition/age/decay, scope by decay — while never saying either value is immutable.

**Q2A-13 (Source-derived — K-R8, `:222`).** "`knowledge_class` is a field on memories **and** assets. Ambient applies
only to authored world/scenario definitions." This constrains **what may be ambient** and ties ambient status to
authorship; it does not fix the value for all time.

**Q2A-14 (Source-derived — I13, `:299`).** "Knowledge class is authored, not inferred from popularity or age." Class
has an author. The rule says nothing about **when** or **how often** authoring occurs; it excludes inference, not
re-authoring.

**Q2A-15 (Source-derived — formation assigns initial values).** §4 Formation assigns "scope, knowledge class,
attribution, time, status" (`:279`); §11 writes `provisional` with class and attribution (`:604`); §8.1 writes
`scope=world`, `knowledge_class=experiential`, `status=active` (`:404`). These describe **initial assignment**; none
says "once" or "thereafter fixed."

**Q2A-16 (Source-derived — the record carries the fields).** §5 lists `knowledge_class` and `scope` as fields
(`:316`, `:318`); §6.1 lists them in `MemoryRecord` (`:340`). A field can hold a later value; nothing in the field
list makes it immutable.

**Q2A-17 (Observed — the same text in a second path).** The normative extract reproduces K-R12 (`:106`) and K-R10
(`:156`) verbatim. That is **wording continuity**, not an independent witness (CLV-07/08).

**Q2A-18 (Observed — no "fixed" text).** A search of r33 and the normative extract for statements that class or scope
are fixed, immutable or assigned once returns **none**. The only prohibitions are cause-specific (repetition, age,
decay, popularity).

---

## 3. The strongest counter-reading

**Q2A-19 (Calculated — the reading that would defeat Yes).** K-R12's "plus `scope` and `knowledge_class` changes"
could be read as a **dimension label** — the stream is authoritative for these dimensions — rather than as an
assertion that changes occur. On that reading the clause is compatible with immutable values, and Q2 would be
**Unestablished**.

**Q2A-20 (Calculated — why that reading is weaker).** (i) The word is **changes**, not "values"; a dimension label
would not need the process noun. (ii) The same sentence contrasts the projection with "mutable authoritative
fields" — the contrast is about **how** values change, not whether they do. (iii) K-R9, K-R10 and D5 independently
regulate the **causes** of class and scope change; on the immutable reading all three are redundant, and K-R10's
"never downgrades … by repetition or age" is a prohibition on a non-event. (iv) Three ratified or semi-ratified rules
converging on the same presupposition is stronger than one ambiguous noun.

**Q2A-21 (Calculated — what the counter-reading still achieves).** It keeps open the **sufficiency** and
**mechanism** questions, which is correct: Q2 = Yes does not say what may cause a change, how it is recorded, or
whether any particular change is warranted. It says only that the model's domain includes such changes.

**Q2A-22 (Observed — the counter-reading is recorded, not dismissed).** If the owner reads K-R12 as a dimension
label, Q2 returns to Unestablished and the carrier analysis stays conditional. The evidence below is presented so
that reading can be tested rather than assumed.

---

## 4. Original design intent

**Q2A-23 (Observed — no revision chronology).** r33 is untracked and has no git history (CLV-01). Its header's
revision log does not mention §5's status row or §11's diagram (CLV-02); **no entry mentions K-R12 or §8.1**
(EKR-26).

**Q2A-24 (Unestablished — why the phrase is there).** There is **no original explanation** of why K-R12 says "plus
`scope` and `knowledge_class` changes." Which revision introduced the phrase is not established (CLV-04 for §11; the
same absence applies to K-R12). The intent behind the wording cannot be recovered from the surviving record.

**Q2A-25 (Inferred — what the surrounding rules imply about intent).** K-R9's "popularity is not authorship," K-R10's
"never downgrades … by repetition or age," and D5's "decay never affects scope" are all **cause-regulating** rules.
Their presence is evidence that the designers contemplated class and scope change and wanted to constrain it, not
that they thought it impossible. This is an inference from rule design, not a statement of intent.

**Q2A-26 (Source-derived — an adjacent, non-applicable precedent).** Amendment 02's `OD-04` (2026-11-11) establishes
**explicit scope promotion** for *conversation state*: a change is preserved at conversation scope and "may enter
character-persistent memory or automatic lore only through the explicit promotion required by `OD-04`; the entry
records its source scope, prior scope, and promotion reason" (`REVISION-8 :523`).

**Q2A-27 (Calculated — why that precedent cannot be transferred).** OD-04 governs **conversation-continuity state**,
not `MEMORY-09` memory records. It is a different object, a different lifecycle and a different authority chain. It
shows that the project has a working concept of explicit scope promotion elsewhere; it is **not** evidence for r33's
memory stream and cannot supply a carrier, a rule or an authority.

**Q2A-28 (Source-derived — a checkpoint boundary, also non-canonical).** The specialised-memory checkpoint states:
"Scenario-specific information must not become world-wide information" (checkpoint §7). This is a boundary statement
in a research checkpoint, not r33 canon; it presupposes that such a broadening is possible and forbids it. It
corroborates contemplation, not permission.

---

## 5. Behavioural examples

**Q2A-29 (Observed — no worked change).** r33's worked cases assign class and scope at formation and never change
them: §8.1's earthquake is formed `world` / `experiential` / `active` (`:404`); §8.2's north gate is
`ambient_authored`; §8.3's haunted house is `restricted_authored` (§8.2–§8.3 `:410`–`:420`). §8.4 changes **status**
(retraction), not class or scope. **No example shows a post-formation class or scope change.**

**Q2A-30 (Observed — no conformance requirement).** The 118-row map contains **no** row requiring or testing a
class/scope change; its scope-related row is L7, which concerns authority grants in a domain (authority-promotion),
not the memory `scope` field.

**Q2A-31 (Calculated — what the absence does and does not prove).** No example and no requirement is **not** evidence
that change cannot occur. It is evidence that the behaviour is **unexemplified and untested** — a representation and
coverage gap, consistent with S-C1. The owner's caution applies: absence of a kind or example is not proof of
immutability.

**Q2A-32 (Source-derived — the one worked scope operation is elsewhere).** OD-04's promotion records source scope,
prior scope and promotion reason (Q2A-26). It is a complete worked operation — in the **conversation** system, not
`MEMORY-09`. Its completeness is a contrast, not a template.

---

## 6. Threshold applied

**Q2A-33 (Calculated — the finding).** **Q2 = Yes.** The ratified text explicitly names `scope` and
`knowledge_class` **changes** as contents of the lifecycle stream (K-R12), and three rules regulate the **causes** of
such changes (K-R9, K-R10, D5). No text fixes either value after formation. The counter-reading (Q2A-19…Q2A-22) is
possible but weaker, and it leaves the cause-regulating rules unexplained.

**Q2A-34 (Calculated — what Yes establishes).** Only **possibility within the model**: the domain includes
post-formation class and scope changes. It does **not** establish that any change is required, how often one occurs,
what suffices to warrant one, who may cause one, or how one is represented.

**Q2A-35 (Observed — what Yes does not touch).** It does not answer C-3, does not decide whether N-a covers
class/scope, does not select a carrier, and does not create a permission, a rule or a requirement.

**Q2A-36 (Calculated — Q2 = Yes is a research finding, not an owner decision).** It is the threshold applied to the
text, and the owner may contest the reading (Q2A-19). It adds nothing to the map and changes no canonical state.

---

## 7. Consequence: the carrier is now required, and none is admissible

**Q2A-37 (Calculated — the conditional resolves).** C3D-18/19 made Q5 conditional on Q2. With Q2 = Yes, the
conditional's first branch is taken: **a positive carrier is required** for class and scope changes.

**Q2A-38 (Calculated — and the current constraint set supplies none).** Under NSD-09 and A-6, P-1 (own kind), P-2
(generic entry) and P-6 (piggyback attribute) are barred; P-4 is strained by "of that stream" and K-R14
independence; P-5 is storage, not a carrier (C3D-05…C3D-11). **No positive carrier is admissible today.**

**Q2A-39 (Calculated — the resulting tension, stated as a finding).** K-R12 provides for class/scope changes, and the
current constraint set cannot carry them. The null carrier (P-3) remains the only **buildable** projection, but it is
now known to be **incomplete relative to K-R12**, not merely minimal. This is a **normative-representational
tension**, not a semantic prohibition: the constraint set, not the meaning of class or scope, is what blocks the
carrier.

**Q2A-40 (Calculated — what the tension does not authorise).** It does not authorise a new kind, field or schema, and
it does not select a carrier. It identifies a conflict the owner would have to resolve if class/scope changes are to
be representable: revise the constraint set, assign the carrier to a part outside N-a, or accept that the stream's
declared domain exceeds what it can record.

**Q2A-41 (Calculated — the effect on C-3).** Q2 = Yes removes the "nothing changes" branch. A **status-only N-a** can
no longer be justified on the ground that class and scope never change; it can only be a **scope boundary** that
leaves a required carrier outside N-a. Q5 is now live rather than conditional, but Q1 (which claim C-3 is about) and
the carrier choice remain open, so **Q4 (deferred or excluded) is still premature** (C3D-21…C3D-26).

**Q2A-42 (Calculated — the effect on the null carrier's status).** "The null carrier is the only currently buildable
representation" remains true (C3D-11). "The null carrier is the complete intended model" is now **false**: K-R12's
domain includes changes it cannot carry. That is exactly the distinction the owner drew before the audit.

---

## 8. What this record does not do

**Q2A-43 (Observed).**

| Untouched | State |
|---|---|
| **C-3** | **not answered**; no scope determination |
| **Q1 / Q4 / Q5** | open; Q5 now required, not conditional; Q4 still premature |
| **carrier** | none selected; no positive carrier admissible today |
| **A-6 / NSD-09** | **not relaxed**; no kind, field, key or schema proposed |
| **C-1 / C-2** | determined / non-decision; not reopened |
| **N-a / F-1…F-6** | unanswered; **F-6 is not a default** |
| **N-b / N-c / entry architecture** | untouched |
| **O-8a / DN-6a / DN-6b / DN-1 / DQ-4** | unchanged |
| **§11's adoption; §8.1's collision** | not decided |
| **map / requirements** | **118**; DQ-3's capture gap acknowledged |

**Q2A-44 (Observed — no rule, predicate, schema or capture).** No authority rule, validity predicate, event kind,
field, key, fold, entry architecture or canonical record is proposed or created.

---

## 9. Traceability

**Q2A-45 (Observed).**

| What | Where |
|---|---|
| K-R12: stream kinds "plus scope and knowledge_class changes"; projection | `:111`; normative extract `:106` |
| K-R10: no downgrade by repetition/age; decay does not change class | `:224`; extract `:156` |
| K-R9: no ambient by popularity; popularity is not authorship | `:223` |
| D5: decay never affects scope | `:72` |
| K-R8: class field; ambient only authored definitions | `:222` |
| I13: class authored, not inferred | `:299` |
| Formation assigns initial values | `:279`, `:604`, `:404` |
| Fields in the record | `:316`, `:318`, `:340` |
| Worked examples assign at formation only | §8.1–§8.4 `:404`–`:430` |
| No revision note for K-R12/§8.1; no git history | CLV-01, CLV-02; EKR-26 |
| Adjacent scope-promotion operation (different system) | Amendment 02 `OD-04`; REVISION-8 `:523` |
| Checkpoint boundary on scenario→world | checkpoint §7 |
| No class/scope change in the map | 118-row map; L7 is authority-scoped |
| Null vs positive carrier; P-5 is storage | C3D-05…C3D-11 |
| Q5 conditional on Q2; order Q1→Q2→Q5→Q4 | C3D-18/19; C3D-21 |
| The owner's threshold and cautions | 2026-11-20; Q2A-01…Q2A-03 |

---

## 10. Limits and record status

**Q2A-46 (Observed — method).** **Newly read:** r33 `:64`–`:81`, `:205`–`:236`, `:275`–`:300`, `:315`–`:345`,
`:404`–`:430`, K-R12/K-R10/K-R9/K-R8/D5/I13; the normative extract's K-R12/K-R10; CLV-01/02; the specialised-memory
checkpoint §7; Amendment 02 `OD-04` and its revision-8 scope-promotion passage; the 118-row map's scope rows. **Not
inspected:** any implementation's class/scope handling, r32's bytes (absent), the frozen preregistration, v9 beyond
§7.

**Q2A-47 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The audit is textual and historical.

**Q2A-48 (Observed — the labels' limits).** Q2A-05, Q2A-07, Q2A-09, Q2A-11, Q2A-13…Q2A-17, Q2A-23, Q2A-26,
Q2A-28…Q2A-30, Q2A-32, Q2A-35, Q2A-43…Q2A-45 are **Observed** or **Source-derived**. Q2A-06, Q2A-08, Q2A-10,
Q2A-12, Q2A-19…Q2A-22, Q2A-27, Q2A-31, Q2A-33…Q2A-34, Q2A-36…Q2A-42 are **Calculated** or **Inferred** except where
marked **Unestablished**. The strongest counter-reading is recorded at Q2A-19…Q2A-22, not suppressed.

**Q2A-49 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Q2 finding | **Yes** — a research finding, not an owner decision |
| Carrier selected; A-6 relaxed | **0** |
| C-3 answered; request drafted | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** **Q2 = Yes** (Q2A-33). K-R12 explicitly names **scope and `knowledge_class` changes** as contents
of the lifecycle stream (`:111`), and K-R9, K-R10 and D5 regulate the **causes** of such changes — rules that would
be redundant if the values were fixed (Q2A-06…Q2A-12). No text fixes either value; the "fixed" search is empty
(Q2A-18). The counter-reading — that "changes" is a dimension label — is possible but weaker, and it leaves the
cause-regulating rules unexplained (Q2A-19…Q2A-22). **The finding establishes possibility only** (Q2A-34). **A
positive carrier is therefore required, and none is admissible under the current constraints** (Q2A-37/38): the null
carrier remains buildable but is now known to be **incomplete relative to K-R12** (Q2A-39/42). **Q5 is live; Q1 and
the carrier remain open; Q4 is premature; C-3 is not answered** (Q2A-41). **Nothing is selected, requested, designed
or authorised.**
