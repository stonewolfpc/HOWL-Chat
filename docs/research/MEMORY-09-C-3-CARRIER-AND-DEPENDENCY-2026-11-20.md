# MEMORY-09 — C-3 Carrier and Dependency

## P-5 is storage, not a carrier — and the order that follows from it

**Status:** question-analysis stage. **Research evidence only — no canonical edit, no schema, no code, no test, no
probe, no ledger entry, no implementation.** **C-3 is not answered; no carrier is selected; no owner decision request
is drafted.** C-1 is not reopened; C-2 stays a non-decision. No fold shape, entry architecture, domain, event kind,
field, key, requirement or schema is selected.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Refinement | **P-5 is a storage form, not a complete carrier explanation**; it tells where a derived value is kept, not what changes it or how the projection is rebuilt |
| Consequence | under the current constraints there is **no positive carrier** for a class/scope change; the only admissible projection is the **null carrier** (formation values + status transitions) |
| Dependency order | **Q1 (which claim) → Q2 (semantic possibility) → Q5 (carrier, conditional) → Q4 (disposition)**; Q3 is a finding that feeds Q5 |
| Does not do | answer C-3; choose a carrier; revise the constraints; draft a request |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. The task, and what is carried forward

**C3D-01 (Observed — the owner's refinement).** P-5 is a **storage form**, not yet a complete carrier explanation: it
tells us where a derived value might be retained, not what authoritative input changes it or how the projection is
rebuilt. It should not be treated as a complete alternative to P-3 without that missing explanation.

**C3D-02 (Observed — the dependency order stated).** Establish whether C-3 concerns **derivation**, **transition
representation**, or **both**; establish what **carrier** K-R12 actually requires or permits; only then determine
whether class/scope are **deferred or excluded** from N-a. Keep **semantic possibility** separate from **current
representability** throughout.

**C3D-03 (Observed — carried forward).** C-1 is determined (R-E, §11-scoped) and not reopened. C-2 is a non-decision.
S-C1 is Observed; S-C2 is Unestablished (C3Q-19…C3Q-22). N-b, N-c and the entry architecture remain open. r33 `scope`
and v9 authority `domain` remain distinct axes. No new kinds, fields, keys, requirements or schema (NSD-09; A-6).

**C3D-04 (Observed — no selection).** No carrier is chosen and the constraints are not revised (C3D-29).

---

## 2. P-5, corrected

**C3D-05 (Source-derived — what a stored projection is).** §6.4's check folds events, then "compare[s] to stored
projection" (`:366`). A2-11 reads §6.1's `status` field charitably as a **cached projection** rather than a second
authority. That is a statement about **storage and comparison**, not about **input**.

**C3D-06 (Calculated — a cache presupposes three things it does not supply).**

| What a cache needs | What P-5 says about it |
|---|---|
| an **authoritative input** that changes the value | nothing |
| a **rebuild rule** that recomputes the value from that input | nothing |
| a **check** that compares computed to stored | names the check (`:366`) but not the computation |

Without the first two, P-5 is a location, not an explanation.

**C3D-07 (Calculated — P-5 is orthogonal to the carrier).** P-5 can accompany **any** carrier: P-1, P-2, P-3 or P-6.
It does not compete with them and cannot substitute for the missing one. The earlier phrasing "P-3 and P-5 are the
only potentially representable readings" (C3P-19) is **Narrowed**: P-5 is not a reading of the carrier at all.

**C3D-08 (Calculated — with a null carrier, the cache is trivial).** If the only projection is formation values plus
status transitions, then a stored class/scope value is simply the formation value. The check at `:366` recomputes it
without consulting any class/scope input, so it cannot detect a **missing** class/scope change: there is no expected
change to compare against. A cache over a null carrier is not evidence of a carrier.

**C3D-09 (Calculated — the correct taxonomy).**

| Layer | Question | Members |
|---|---|---|
| **carrier** | what authoritative input changes the value? | **null** (P-3) or **positive** (P-1, P-2, P-4, P-6) |
| **storage** | where is the derived value retained? | **P-5** (cache), or none |
| **check** | how is the stored value compared to the computed one? | `:366` |

The three layers are distinct. P-5 belongs to the second; it cannot be counted as a member of the first.

---

## 3. The revised carrier space

**C3D-10 (Calculated — null versus positive).**

| Class | Reading | What it says about a class/scope change | Status under current constraints |
|---|---|---|---|
| **null** | **P-3** — formation values only | no post-formation change is carried | **admissible** (nothing new needed) |
| **positive** | **P-1** — own transition kind | a new lifecycle kind carries it | barred — new kind |
| **positive** | **P-2** — generic change entry | a general change entry carries it | barred — new kind |
| **positive** | **P-4** — another operation joined at projection | an external operation carries it | strained — "of that stream" and K-R14 independence |
| **positive** | **P-6** — piggyback on an existing entry | an existing transition carries the value | barred now — new field |

**C3D-11 (Calculated — the corrected count).** Under the current constraints there is **no admissible positive
carrier** for a class/scope change. The only admissible carrier is the **null carrier**, and P-5 is a storage option
over it. The earlier "P-3/P-5" formulation counted a storage layer as a carrier; corrected, the space is **one null
carrier plus optional storage**.

**C3D-12 (Narrowed — what that does not mean).** It does **not** mean class and scope are semantically immutable
(S-C2 stays Unestablished), and it does **not** answer C-3. It means: **no positive carrier is buildable today without
a constraint revision**. That is a representability finding, and the constraints are the owner's to revise
(C3P-20).

**C3D-13 (Calculated — P-4's strain is independent of the schema bar).** P-4 fails on the architecture's own terms
before any constraint is considered: K-R14 rebuilds each stream independently (`:113`), `:367` isolates the
projections, and K-R13's cross-stream use is a **read-time join** of edge × memory projection, not a second input to
the memory projection (`:112`). A class/scope change recorded elsewhere would either not be rebuildable from the
lifecycle stream or would have to be written into it — collapsing P-4 into P-1/P-2/P-6.

**C3D-14 (Calculated — P-6 is the nearest positive carrier).** Piggybacking needs no new **kind**; it needs an
attribute on an existing entry. That is why it is the positive reading closest to admissibility, and why the field
constraint, not the kind constraint, is what blocks it. Naming this is not a recommendation; it identifies which
constraint would have to change first if P-6 were ever considered.

**C3D-15 (Calculated — no positive carrier is "more faithful").** The readings differ in what they make
representable, not in fidelity to K-R12. Claim A is satisfiable by the null carrier; claim B is not (C3D-17).

---

## 4. What claim A requires, and what it does not

**C3D-16 (Source-derived — claim A).** K-R12 requires current status, class and scope to be **derived**, not mutable
authoritative fields (`:111`). That is the projection requirement.

**C3D-17 (Calculated — derivation is compatible with a null carrier).** If class and scope never change after
formation, their current value is still **derived**: it is the formation value, recomputed rather than treated as a
mutable authority. Claim A is satisfied trivially. **Claim A does not require a carrier for every dimension**; it
requires that whatever the value is, it is not a mutable authoritative field.

**C3D-18 (Calculated — but claim A does require a carrier if the value changes).** If a class or scope value can
change, the new value must come from somewhere authoritative for the projection. Otherwise the change would either be
a mutable field write (barred by claim A) or an unrepresented difference (which the projection cannot derive). So:

| Q2 — can it change? | Carrier needed? | Current representation |
|---|---|---|
| **no** | no | null carrier is complete |
| **yes** | **yes** | none admissible — a positive carrier is required and barred |
| **Unestablished** | **conditional** | null carrier is the only buildable one; its completeness is Unestablished |

**C3D-19 (Calculated — this is the real dependency).** Q5 (carrier) is **conditional on Q2** (semantic possibility).
It is not a peer question. If Q2 is "no," there is nothing to carry; if Q2 is "yes," a carrier is mandatory and the
current constraint set supplies none; if Q2 is Unestablished, the carrier question is conditional and cannot be
answered without deciding Q2 or revising the constraints.

**C3D-20 (Calculated — and this is why P-5 cannot substitute).** A cache does not answer the conditional. It records a
result whichever way Q2 goes; it cannot make a change representable, and it cannot make an unrepresented change
visible.

---

## 5. The revised dependency order

**C3D-21 (Calculated — the owner's four steps mapped to the questions).**

| Step | Question | What it fixes | Current state |
|---|---|---|---|
| **1** | **Q1** — does C-3 concern **derivation** (A), **transition representation** (B), or **both**? | what the answer is an answer to | unchosen |
| **2** | **Q2** — can class/scope change after formation? | whether a carrier is needed at all | **Unestablished** |
| **3** | **Q5** — which carrier does K-R12 require or permit? | the positive carrier, if one is needed | **conditional on step 2**; none admissible today |
| **4** | **Q4** — deferred or excluded from N-a? | the disposition | only meaningful after 1–3 |

**C3D-22 (Calculated — Q3's place).** Q3 (current representation) is not a separate step in the order; it is the
**finding** that feeds step 3: there is no positive carrier today, and the null carrier is the only buildable one
(C3D-11). It does not decide Q2 and does not decide Q4.

**C3D-23 (Calculated — why Q1 must precede Q2).** If C-3 is about **derivation** (A), then Q2 is not part of it: the
question is whether class and scope are derived at all, and the null carrier answers it. If C-3 is about **transition
representation** (B), Q2 is constitutive: a transition requires something that changes. If C-3 is about **both**, the
two answers must be stated separately. The same "status-only" sentence would mean different things in each case.

**C3D-24 (Calculated — why Q2 must precede Q5).** Without Q2, "no carrier" and "carrier not yet chosen" are
indistinguishable. Q2 = no makes the null carrier complete; Q2 = yes makes it incomplete. The dependency is
one-directional: Q5 cannot answer Q2, and the current representation cannot answer Q2 either (C3Q-16).

**C3D-25 (Calculated — why Q5 must precede Q4).** "Deferred" and "excluded" are dispositions of a subject and a
carrier. If the carrier is null, excluding class/scope from N-a is coherent and says nothing about change elsewhere;
if a positive carrier exists, exclusion from N-a leaves the change in the stream but outside the fold; if the carrier
is unestablished, the disposition inherits that uncertainty (C3Q-31; C3P-25).

**C3D-26 (Calculated — the order is not a decision tree yet).** It is a sequence of questions, each of which can be
answered "Unestablished." An answer of Unestablished at step 2 does not force a choice at step 3; it makes step 3
conditional and step 4 premature. That is the correct state, not a defect.

---

## 6. What this record does not do

**C3D-27 (Observed).**

| Untouched | State |
|---|---|
| **C-3** | **not answered**; no carrier chosen; no request drafted |
| **Q1…Q5; T-1…T-6** | open; the order is refined, the answers are not |
| **P-1…P-6** | none selected, ranked or excluded; P-5 reclassified as storage |
| **S-C2** | **Unestablished**; the null carrier does not make it true |
| **C-1 / C-2** | determined / non-decision; not reopened |
| **N-a / F-1…F-6** | unanswered; **F-6 is not a default** |
| **I-1 / I-2 / I-3** | I-3 excluded by K-R12; the others open |
| **N-b / N-c / N-d** | untouched |
| **O-8a / O-9 / O-10b / O-11** | unchanged; authority-promotion only |
| **DN-6a / DN-6b; DN-1; DQ-4** | unchanged / settled / open |
| **constraints** | **not revised**; no kind, field, key or schema proposed |
| **map / requirements** | **118**; DQ-3's capture gap acknowledged |

**C3D-28 (Observed — no rule, predicate, schema or capture).** No authority rule, validity predicate, event kind,
field, key, fold, entry architecture or canonical record is proposed or created.

**C3D-29 (Observed — no selection and no revision).** The record does not select a carrier, does not revise NSD-09 or
A-6, and does not recommend that either be revised. It states what would be needed if a positive carrier were ever
considered.

---

## 7. Traceability

**C3D-30 (Observed).**

| What | Where |
|---|---|
| The projection clause; claim A | K-R12 (`:111`) |
| Claim B is not entailed | C3P-06/07; C3P-08 |
| The six readings | C3P-10 |
| The earlier "P-3/P-5" count, now narrowed | C3P-19; C3D-07/11 |
| Stored projection and the check | `:366`; A2-11 |
| Independent rebuild; projection isolation; read-time join | K-R14 (`:113`); `:367`; K-R13 (`:112`) |
| No cross-stream status conferral | ECI-16 |
| Formation writes class/scope | `:604`; `:404`; ECI-04/05 |
| S-C1 / S-C2 | C3Q-19…C3Q-22 |
| Q1 and Q5 | C3Q-14; C3P-23/24 |
| Q5 prior to Q4 | C3Q-31; C3P-25 |
| No new kinds/fields/requirements | NSD-09; A-6 (NAE-31) |
| Scope vs authority domain | RX-17; C3Q-05 |
| The owner's refinement and order | 2026-11-20; C3D-01/02 |

---

## 8. Limits and record status

**C3D-31 (Observed — method).** **Newly read:** K-R12 (`:111`), K-R13 (`:112`), K-R14 (`:113`), `:366`/`:367`, §6.1
(`:341`), §11 (`:604`), §8.1 (`:404`), A2-11/12, ECI-04/05/16, C3P-10/18/19/23…26, C3Q-14/16/31, NSD-09, A-6. **Not
inspected:** any implementation's class/scope handling or projection rebuild, r32's bytes, the frozen preregistration,
v9 beyond §7. No test, probe or model call was run (C3D-32).

**C3D-32 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The record is textual and definitional.

**C3D-33 (Observed — the labels' limits).** C3D-01…C3D-05, C3D-16, C3D-27…C3D-30 are **Observed** or
**Source-derived**. C3D-06…C3D-15, C3D-17…C3D-26 are **Calculated** except where marked **Narrowed** or
**Unestablished**. Nothing is offered as a provision of r33 beyond the quotations already recorded.

**C3D-34 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| C-3 answered; carrier chosen; request drafted | **0** |
| Constraints revised; fold, entry model or domain selected | **0** |
| C-1 reopened; C-2 made a decision | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** **P-5 is a storage form, not a carrier**: it names where a derived value is retained, not what
changes it or how the projection is rebuilt (C3D-05…C3D-09). **The carrier space is therefore null (P-3) versus
positive (P-1/P-2/P-4/P-6)**, and under the current constraints **no positive carrier is admissible** — the only
buildable projection is the null carrier plus optional storage (C3D-10…C3D-12). **Claim A is compatible with the null
carrier**, because derivation does not require a carrier for a dimension that never changes (C3D-16/17); **claim A
does require one if the value changes**, so **Q5 is conditional on Q2** (C3D-18/19). **The order is Q1 → Q2 → Q5 →
Q4**, with Q3 as a finding that feeds Q5, and "Unestablished" is a permissible answer at every step (C3D-21…C3D-26).
**Nothing is answered, chosen, requested or authorised; the constraints are not revised.**
