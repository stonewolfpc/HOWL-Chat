# MEMORY-09 — C-3 Question Analysis

## What "is N-a status-only?" is actually asking

**Status:** question analysis. **Research evidence only — no canonical edit, no schema, no code, no test, no probe,
no ledger entry, no implementation.** **C-3 is not answered and no owner decision request is drafted.** C-1 is taken
as determined and is not reopened. **C-2 is dispositioned as a non-decision** (consolidation/form only). No fold
shape, entry architecture, domain, event kind, field, key, requirement or schema is selected.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Task | establish **exactly what C-3 asks**, before any answer is chosen |
| Finding | "status-only" bundles **five distinct questions** — subject, sequencing, semantics, representation, and operation identity — and they do not share an answer |
| Correction | the missing class/scope vocabulary is a **current representation gap**, not proof that class and scope cannot change (**Narrowed**) |
| Does not do | answer C-3; choose S-A/S-B/S-C; draft a decision request; select a fold or entry model |
| Preserved | N-b, N-c, entry architecture open · r33 `scope` ≠ v9 authority `domain` · no new kinds/fields/keys/requirements/schema |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. The task, and what is carried forward

**C3Q-01 (Observed — the owner's direction).** The next useful work is to establish **what C-3 is asking**, not to
choose an answer; the point is to pause before drafting an owner decision request.

**C3Q-02 (Observed — C-1).** §11's "authority" is epistemic support for record-promotion, §11-scoped (C1D-01…C1D-07).
The four Unestablished items stand, including permission to perform record-promotion (C1D-13). C-1 is not reopened.

**C3Q-03 (Observed — C-2).** C-2 is dispositioned as a **non-decision**: every substantive constraint is already
binding, so a consolidation or note is form only and creates no authority (C23-05…C23-10; owner agreement
2026-11-20). No consolidation artefact is created by this record.

**C3Q-04 (Observed — what remains open).** **N-b** (`established`↔`active`; `provisional` representation), **N-c**
(status-stream ordering), and the **entry architecture** (I-1 / I-2; I-3 excluded by K-R12) remain open. Neither C-2
nor C-3 determines them (C23-22/23, C23-24).

**C3Q-05 (Observed — the two axes that must not merge).** r33's `scope` field (`world` / `scenario` /
`character_public` / `character_private`, §5 `:318`) is not v9's authority **domain** (the per-world/scenario domain
of O-8a). RX-17 records them as different axes. C-3 is about the former, never the latter (C23-15).

**C3Q-06 (Observed — the standing constraint).** No new event kinds, fields, keys, requirements or schema
(NSD-09; A-6/NAE-31). The map remains at **118**; DQ-3's representation requirement remains the acknowledged capture
gap.

---

## 2. The correction: representation gap is not semantic impossibility

**C3Q-07 (Source-derived — what is actually absent).** K-R12 names the lifecycle stream as
`established | superseded | corrected | retracted`, "plus `scope` and `knowledge_class` changes" (`:111`). The named
kinds are status-shaped. **No kind names a class or scope change**, and no schema exists for any lifecycle entry
(A2-11/A2-12; ECI-11).

**C3Q-08 (Narrowed — the readiness record's inference).** `NAB-26` and `C23-14` said class/scope folds "require new
kinds, which A-6 excludes." That is correct only as a statement about **what can be expressed under the current
constraint set**. It is **not** a finding that class or scope cannot change, and not a finding that a class/scope
change is semantically impossible. The correct statements are:

| Claim | Status |
|---|---|
| a class/scope fold is **not expressible with the current vocabulary** | **Observed / Calculated** |
| creating the vocabulary **now** is excluded by A-6 and NSD-09 | **Observed** |
| class and scope **cannot change** after formation | **Unestablished** — and not implied by the absence |
| class/scope changes are **not intended** to be representable | **Unestablished** — K-R12 says they are stream inputs |

**C3Q-09 (Calculated — why the distinction matters).** If the absence were read as a permanent prohibition, C-3
would be answered by the schema gap: "status-only, because nothing else can exist." That would convert a current
artefact into a semantic conclusion — the move I15 and I18 warn against in other forms, and the move the owner
identified. **The gap establishes a present asymmetry, not a boundary of meaning.**

**C3Q-10 (Calculated — the most C-3 can inherit from the gap).** The gap makes one answer **cheaper** (status-only is
expressible today) and another **costlier** (a class/scope fold needs vocabulary that does not exist). Cost is not
correctness, and C-3 is not a cost question.

---

## 3. What C-3 is about: three candidate subjects

**C3Q-11 (Calculated — the question has no fixed subject yet).**

| Subject | The question it would ask | Why it is distinct |
|---|---|---|
| **N-a** | does the **fold** cover class and scope, or only status? | a scope question about one part |
| **the lifecycle stream** | does the stream **contain** class/scope change entries at all? | a vocabulary question about the whole stream |
| **K-R12's projection requirement** | are class and scope **projected** from the stream, as K-R12 says, and if so from what? | a consistency question about the requirement itself |

**C3Q-12 (Source-derived — the three subjects are not interchangeable).** K-R12 says status, class and scope are "a
projection of that stream" (`:111`). That sentence is wider than N-a: it is a claim about the stream. §6.1 lists
`status`, `scope` and `knowledge_class` as record fields (`:341`), and A2-11 records the contradiction: the
projection claim and the field list cannot both be literally true, and `supersedes`/`superseded_by`/`correction_of`
have no event to project from. So a question about N-a's scope is not yet a question about what the stream contains.

**C3Q-13 (Calculated — C-3's subject must be chosen before its answer means anything).** "Status-only" answered
about N-a leaves the stream question open; answered about the stream, it denies K-R12's three-dimension sentence;
answered about the projection requirement, it becomes a consistency finding rather than a scope decision.

---

## 4. The five questions bundled inside "status-only"

**C3Q-14 (Calculated — the decomposition).**

| # | Question | What an answer would settle | Current state |
|---|---|---|---|
| **Q1 — subject** | is C-3 about **N-a**, the **stream**, or the **projection requirement**? | what the answer is an answer *to* | unchosen (§3) |
| **Q2 — semantics** | **can** class and scope change after formation? | whether there is anything for a fold to compute | **Unestablished**; K-R12 implies yes, no mechanism says how |
| **Q3 — representation** | is there a **current vocabulary** for such changes? | whether the change can be recorded today | **no** — no kind, no schema (C3Q-07) |
| **Q4 — sequencing** | if not in N-a now, are class/scope **deferred** as separate parts or **excluded** from ND-1? | S-B vs S-A | unchosen |
| **Q5 — operation identity** | is a class/scope change a **lifecycle entry**, a **record edit**, or a **different operation**? | which part owns it at all | **Unestablished**; no candidate is clean (§6) |

**C3Q-15 (Calculated — the questions are ordered, not parallel).** Q1 fixes the subject; Q2 and Q3 are factual
questions about the dimension; Q5 fixes the operation; Q4 is a disposition about N-a. Q4 cannot be answered before
Q1 and Q5, because "deferred" and "excluded" are dispositions of a subject that must first be identified.

**C3Q-16 (Calculated — Q2 and Q3 are independent).** A dimension can be changeable in principle and unrepresentable
today (Q2 yes, Q3 no), or representable and semantically fixed (Q2 no, Q3 yes). The current evidence supports
**Q2 Unestablished, Q3 no**. Treating Q3's "no" as Q2's answer is the conflation the owner warned against.

**C3Q-17 (Unestablished — Q2's evidence is thin in both directions).** For changeability: K-R12 names class and scope
changes as stream inputs, and formation writes initial values (`:604`; `:404`), which implies later values can differ.
Against: no operation, kind, or rule describes such a change; §5 lists the values as fields; A2-11 shows the field
form is already in tension with the projection claim. Neither side is established.

**C3Q-18 (Observed — Q3 is a clean absence).** No lifecycle kind for class/scope; no schema; no fold rule; no
projection input. The absence is recorded, not inferred (NAB-23/24; A2-12).

---

## 5. What the earlier S-A / S-B / S-C labels were doing

**C3Q-19 (Observed — the labels).** C23-12 offered **S-A exclusive**, **S-B deferred**, **S-C non-fold** as three
meanings of "status-only."

**C3Q-20 (Calculated — S-A and S-B are dispositions; S-C is not one question).** S-A and S-B are answers to **Q4**
(defer or exclude) once Q1 is fixed. **S-C bundles Q2 and Q3**: "the stream carries only status transitions" can mean
"today's representation carries only status" (Q3) or "class and scope are not lifecycle dimensions at all" (Q2). The
bundle is the source of the overreach in C3Q-08.

**C3Q-21 (Calculated — the corrected mapping).**

| Label | Actually answers | Must be split into |
|---|---|---|
| **S-A** | Q4: excluded from ND-1 | stays |
| **S-B** | Q4: deferred to separate parts | stays |
| **S-C** | Q2 **and** Q3 together | **S-C1** current representation carries only status (Q3) · **S-C2** class/scope are not lifecycle dimensions in principle (Q2) |

**C3Q-22 (Narrowed — S-C's status).** S-C is **Narrowed** to S-C1 (Observed) and S-C2 (Unestablished). Only S-C1 is
supported by the absence of kinds; S-C2 is a semantic claim the absence cannot establish.

---

## 6. Q5 — operation identity, and why it is the hardest

**C3Q-23 (Calculated — candidate homes for a class/scope change, none clean).**

| Candidate home | What it would make the change | Why it does not fit cleanly |
|---|---|---|
| **lifecycle stream** (K-R12) | an entry in the same stream as status transitions | no kind, no schema; the stream's authority is declared but unspecified (A2-12) |
| **record field edit** (§6.1) | a mutable field write | K-R12 forbids mutable authoritative fields for class and scope |
| **authority-promotion** (O-8a) | a separate world/scenario record with declared applicability | it creates a record; it does not change an existing memory's class or scope (O-11; CPA-07) |
| **Q7 ambient authorship** (O-8a-5) | authoring ambient content | it is an authorship operation, not a lifecycle change; the operations are distinct |
| **disclosure** (O-3) | a recipient-specific account | it creates an account for a recipient; it does not change the source's class or scope |

**C3Q-24 (Calculated — so Q5 is Unestablished, not merely unrecorded).** No candidate is both authorised and
well-formed. The lifecycle stream is the only one K-R12 names, and it has no vocabulary; the field route is excluded
by K-R12; the other three are different operations with different objects.

**C3Q-25 (Calculated — why Q5 cannot be skipped).** If a class/scope change has **no operation**, then "N-a covers
class and scope" would require inventing one (barred by A-6), while "N-a excludes them" is true but says nothing
about whether they can change elsewhere. Either way, a bare "status-only" answer would be silent on the question that
makes the domain interesting.

**C3Q-26 (Calculated — Q5 is not an authority question).** Identifying an operation is not deciding who may perform
it. C-1's U-2 (permission for record-promotion) stays open, and no operation identified here would supply it.

---

## 7. What C-3 is not asking

**C3Q-27 (Observed).**

| Not C-3 | Why |
|---|---|
| **who may change** class or scope | that is an authority question; C-1's U-2 is open and O-8a is confined to authority-promotion |
| **how** a class/scope fold computes | that is N-a's shape (F-1…F-6), not its domain |
| whether `scope` and `knowledge_class` **exist as fields** | they do (`:341`); the question is whether they are projections |
| whether a v9 **authority domain** change equals a scope change | they are different axes (RX-17) |
| whether **class/scope values are initialised at formation** | they are (`:604`; `:404`; ECI-04/05) |
| whether the lifecycle stream **should** gain kinds | that is a schema question, barred by A-6 |
| whether **§11** is operative | C-1's U-4; independent of C-3 |

**C3Q-28 (Calculated — the negative boundary protects the question).** Each exclusion removes a way of answering C-3
with material from another question. The most tempting is "who may change class/scope," which would reimport the
authority model C-1 removed.

---

## 8. When C-3 becomes answerable

**C3Q-29 (Calculated — the conditions, stated as terms).**

| # | Term | Why it must precede an answer |
|---|---|---|
| **T-1** | the **subject** — N-a, the stream, or the projection requirement | Q1; the same words answer different questions otherwise |
| **T-2** | **Q2 separated from Q3** — semantic changeability vs current representation | the absence of kinds answers only Q3 |
| **T-3** | the **disposition** — deferred or excluded (S-A / S-B) | Q4; only meaningful after T-1 |
| **T-4** | the **operation** — which part owns a class/scope change, if any | Q5; otherwise "covers" has no referent |
| **T-5** | an explicit statement that **no kind, field, key, requirement or schema** is created | otherwise the answer needs what it excludes |
| **T-6** | what remains **open** — class/scope semantics, N-b, N-c, the entry architecture | stated, not implied |

**C3Q-30 (Calculated — T-2 is the owner's caution in operational form).** Without T-2, an answer can slide from
"there is no vocabulary" to "there can be no change." With T-2, the first is recorded as an observation and the second
must be argued or left Unestablished.

**C3Q-31 (Calculated — T-4 may be prior to T-3).** If a class/scope change has no operation, "deferred" is a
disposition of nothing. If it has an operation elsewhere, "excluded from N-a" is coherent and does not deny the
change. So Q5 can change what Q4 means.

**C3Q-32 (Calculated — C-3 is not ready, and the reason is now specific).** It is not ready because **Q1, Q2/Q3 and
Q5 are unresolved**, not because the owner has not chosen S-A or S-B. A decision request drafted now would force a
disposition before its subject and operation are fixed.

---

## 9. What this record does not do

**C3Q-33 (Observed).**

| Untouched | State |
|---|---|
| **C-3** | **not answered**; no S-A/S-B/S-C chosen; no request drafted |
| **C-2** | dispositioned as a non-decision; no artefact created |
| **C-1** | determined; not reopened |
| **N-a / F-1…F-6** | unanswered; none selected, ranked or excluded; **F-6 is not a default** |
| **I-1 / I-2 / I-3** | I-3 excluded by K-R12; the others open |
| **N-b / N-c / N-d** | untouched |
| **O-8a / O-9 / O-10b / O-11** | unchanged; authority-promotion only |
| **DN-6a / DN-6b; DN-1; DQ-4** | unchanged / settled / open |
| **§11's adoption; §8.1's collision** | not decided |
| **map / requirements** | **118**; DQ-3's capture gap acknowledged |

**C3Q-34 (Observed — no rule, predicate, schema or capture).** No authority rule, validity predicate, event kind,
field, key, fold, entry architecture or canonical record is proposed or created.

---

## 10. Traceability

**C3Q-35 (Observed).**

| What | Where |
|---|---|
| C-1 determination and the four Unestablished items | C1D-01…C1D-13 |
| C-2 as a non-decision | C23-05…C23-10 |
| The three S-meanings | C23-12; C23-19 |
| The earlier "excluded by A-6" inference, now narrowed | NAB-26; C23-14; C3Q-08 |
| K-R12's three dimensions and the kind list | `:111` |
| `scope` values; `status` values; record fields | `:318`; `:322`; `:341` |
| Formation writes class and scope | `:604`; `:404`; ECI-04/05 |
| Projection vs fields; relations without events | A2-11; A2-12 |
| No lifecycle schema | ECI-11 |
| Scope vs authority domain | RX-17; C23-15 |
| Formation is a write, not an event | ECI-11 |
| N-b's vocabulary ownership; `provisional` unresolved | NSD-06/07/18 |
| N-c's dependence on N-a's answer | NSD-12 |
| No new kinds/requirements | NSD-09; A-6 (NAE-31) |
| O-8a's operation and object | O-11; CPA-05…CPA-07 |

---

## 11. Limits and record status

**C3Q-36 (Observed — method).** **Newly read:** C23-12/14/19, NAB-26/28, A2-11/12, RX-17, K-R12 (`:111`), §5
(`:318`), §6.1 (`:341`), §11 (`:604`), §8.1 (`:404`), ECI-04/05/11, NSD-06/07/09/12/18, O-11. **Not inspected:** any
implementation's status or class handling, r32's bytes, the frozen preregistration, v9 beyond §7. No test, probe or
model call was run (C3Q-37).

**C3Q-37 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The record is textual and definitional.

**C3Q-38 (Observed — the labels' limits).** C3Q-02…C3Q-07, C3Q-18, C3Q-22, C3Q-27, C3Q-33 and C3Q-35 are **Observed**
or **Source-derived**. C3Q-08, C3Q-09, C3Q-10, C3Q-11, C3Q-13…C3Q-16, C3Q-20/21, C3Q-23…C3Q-26, C3Q-29…C3Q-32 are
**Calculated** except where marked **Narrowed** or **Unestablished**. Nothing is offered as a provision of r33 beyond
the quotations already recorded.

**C3Q-39 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| C-3 answered; request drafted | **0** |
| Fold, entry model, domain, event kind or schema selected | **0** |
| C-1 reopened | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** **C-3 is not a binary, and it is not ready to be put.** "Status-only" bundles **five questions**:
the **subject** (N-a, the stream, or the projection requirement), **semantics** (can class/scope change?), **current
representation** (is there a vocabulary?), **sequencing** (deferred or excluded), and **operation identity** (which
part owns such a change, if any) (C3Q-14, C3Q-23). The earlier **S-C** bundled semantics with representation and is
**Narrowed** to **S-C1** (observed gap) and **S-C2** (unestablished semantic claim) (C3Q-19…C3Q-22). **The absence of
class/scope kinds is a present representation gap, not proof that class and scope cannot change** (C3Q-08/09).
**Q1, Q2/Q3 and Q5 are unresolved, and T-4 may change what Q4 means** (C3Q-29…C3Q-32). **Nothing is answered,
requested, designed or authorised.**
