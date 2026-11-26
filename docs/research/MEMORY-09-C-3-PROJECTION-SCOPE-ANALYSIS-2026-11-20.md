# MEMORY-09 — C-3 Projection-Scope Analysis

## What "a projection of that stream" entails — and what it does not

**Status:** question-analysis stage. **Research evidence only — no canonical edit, no schema, no code, no test, no
probe, no ledger entry, no implementation.** **C-3 is not answered and no owner decision request is drafted.** C-1 is
not reopened; C-2 stays a non-decision. No fold shape, entry architecture, domain, event kind, field, key,
requirement or schema is selected.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Task | refine C-3 by separating **the projection requirement** from **a transition-per-change requirement** |
| Owner's distinction | K-R12's projection clause does **not** by itself prove that every change to every projected value must be represented as a lifecycle transition |
| Finding | the clause entails **derivation from the stream**, not **one transition kind per changed value**; six readings are distinguishable, and the current constraint set narrows which are representable |
| Does not do | answer C-3; choose a reading; decide Q1/Q5; draft a request |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. The task, and what is carried forward

**C3P-01 (Observed — the owner's refinement).** K-R12's projection requirement does not, by itself, prove that every
change to every projected value must be represented as a lifecycle transition. That is what Q1 (subject) and Q5
(operation identity) need to clarify. C-3 stays a question-analysis stage until T-1…T-6 can be answered without
assuming an operation the evidence has not established.

**C3P-02 (Observed — C-1 and C-2).** C-1 remains determined (R-E, §11-scoped) and is not reopened. C-2 remains
dispositioned as a **non-decision** (consolidation/form only); nothing here turns it into a decision.

**C3P-03 (Observed — the open items).** C-3's five questions remain open: subject, semantics, representation,
sequencing, operation identity (C3Q-14). **N-b**, **N-c** and the **entry architecture** remain open (C3Q-04). r33
`scope` and v9 authority `domain` remain distinct axes (C3Q-05). No new kinds, fields, keys, requirements or schema
(C3Q-06).

**C3P-04 (Observed — S-C's narrowing stands).** S-C1 (current representation carries only status transitions) is
**Observed**; S-C2 (class/scope are not lifecycle dimensions in principle) is **Unestablished** (C3Q-19…C3Q-22). This
record works inside that split rather than reopening it.

---

## 2. The clause, and the two claims it is often asked to carry

**C3P-05 (Source-derived — K-R12, `:111`).**

> "Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus
> `scope` and `knowledge_class` changes — and a memory's current status, class and scope are a **projection of that
> stream**, not mutable authoritative fields."

**C3P-06 (Calculated — claim A, entailed).** The current values are **derived**, not stored as mutable authoritative
fields. The derivation is the projection; the fields in §6.1 are at most a cached result (A2-11's charitable
reading).

**C3P-07 (Calculated — claim B, not entailed).** That **every change to every projected value** is represented as a
**lifecycle transition** — an entry of its own, of a transition kind, in the same enumeration as `established`,
`superseded`, `corrected`, `retracted`. The clause says the values are a projection **of the stream**; it does not say
the stream's only entries are transitions, nor that each changed value needs its own kind.

**C3P-08 (Source-derived — the wording distinguishes them).** The sentence itself separates the **kind list**
(`established | superseded | corrected | retracted`) from "**plus** `scope` and `knowledge_class` **changes**." The
word is *changes*, not *transitions*, and the phrase is additive. That is textual support for the distinction the
owner drew — not a resolution of what those changes are.

**C3P-09 (Calculated — why the conflation was tempting).** If "projection of that stream" is read as "everything
projected is produced by a transition entry," then class/scope changes must have transition kinds. That reading makes
the absence of kinds look like a semantic prohibition. The owner's distinction blocks that step: the absence is a
**representation** fact (S-C1), not a **semantic** one (S-C2).

---

## 3. Six readings of the projection clause

**C3P-10 (Calculated — the readings, stated as properties).**

| # | Reading | How a class/scope change would be carried | Compatible with "of that stream"? |
|---|---|---|---|
| **P-1** | **own transition kind** per change | a new lifecycle kind, e.g. a class-change entry | yes |
| **P-2** | **generic change entry** | one general "field changed" entry naming the dimension and value | yes |
| **P-3** | **static after formation** | it would not be carried; the projection uses formation values only | yes, trivially |
| **P-4** | **multi-source join** | recorded by another operation and joined at projection time | **strained** — "of that stream" says the source is this stream |
| **P-5** | **stored cache** | the record field is a cache of the computed projection | yes — a storage detail, not a source claim |
| **P-6** | **piggyback on existing entries** | a class/scope value rides as an attribute of an existing transition (e.g. a correction) | yes, if the entry carries it |

**C3P-11 (Calculated — P-1 is the strong reading and is not entailed).** P-1 is the reading that would make the
missing kinds decisive. C3P-06/07 show the clause does not require it. P-1 remains a possible architecture; it is not
the clause's only meaning.

**C3P-12 (Calculated — P-3 is the minimal reading and is entailed by nothing).** P-3 says class and scope never change
after formation. K-R12's "plus … changes" points against it as a description of the intended stream, but P-3 is the
only reading that needs no vocabulary beyond what exists. It is a **representational** minimum, not a semantic finding.

**C3P-13 (Calculated — P-4 is the reading the architecture's independence rules resist).** K-R14: "Each stream is
rebuilt and checked independently" (`:113`); `:367`: "Neither projection may consult the other." A memory projection
that reads another operation's log would not be independently rebuildable from its own stream. K-R13's cross-stream
join is a **read-time join** of edge × memory projection, not a second input to the memory projection. So P-4 would
need the other operation to write into the lifecycle stream — which collapses it into P-1, P-2 or P-6.

**C3P-14 (Calculated — P-5 is not a rival).** The stored-field question is about **where the result lives**, not
**what produces it**. A2-11 already reads §6.1's field as a cached projection compared by the check at `:366`. P-5
is compatible with P-1, P-2, P-3 or P-6; it does not choose among them.

**C3P-15 (Calculated — P-6 is the reading that preserves "of that stream" without a new kind).** If an existing
transition entry can carry a class or scope value, the change is represented and projected without a new kind. It
requires the entry to have somewhere to put the value — an event attribute — which is a **field** question, not a
kind question.

**C3P-16 (Calculated — the readings are not exhaustive of design, only of the clause's grammar).** They are the ways
the sentence can be made true. A design could combine them (e.g. P-6 for class, P-2 for scope). Nothing here selects a
combination.

---

## 4. What the current constraint set admits

**C3P-17 (Observed — the constraints).** No new event kind, field, key, requirement or schema (NSD-09; A-6/NAE-31);
no cross-stream status conferral (ECI-16); each stream independently rebuildable (K-R14); projections isolated
(`:367`). The map stays at 118.

**C3P-18 (Calculated — admissibility, per reading).**

| Reading | What it needs | Status under the current constraint set |
|---|---|---|
| **P-1** | a new lifecycle kind | **not admissible** — new kind barred |
| **P-2** | a generic change entry | **not admissible** — new kind barred |
| **P-3** | nothing beyond formation values and status kinds | **admissible**; class/scope static post-formation |
| **P-4** | another operation as a projection input | **not admissible** — conflicts with "of that stream" and K-R14 independence |
| **P-5** | a stored cache | **admissible** as a storage detail; source still needs a reading |
| **P-6** | an attribute on an existing entry | **not admissible now** — new field barred; no new kind needed |

**C3P-19 (Calculated — the resulting narrow space, stated carefully).** Under the **current** constraint set, the
representable readings are **P-3** and **P-5**, and P-5 only as storage. That means class/scope would be **static
after formation** in any representation buildable today. This is a **constraint finding**, not a semantic conclusion:
it says what can be built without new vocabulary, not what class and scope are.

**C3P-20 (Narrowed — and it is not a prohibition).** The constraints are the owner's to revise. A future authorisation
of a kind or field would reopen P-1, P-2 or P-6. **C3P-19 must not be recorded as "class/scope cannot change,"** and
S-C2 remains Unestablished (C3Q-08/09).

**C3P-21 (Calculated — P-3 is not the same as S-C2).** P-3 is a statement about a **representation** (no post-formation
change is carried). S-C2 is a statement about **lifecycle dimensions** (class/scope are not lifecycle properties at
all). The first can be true while the second is false: class/scope could be changeable in principle and merely
unrepresented today.

**C3P-22 (Calculated — and P-3 is not an answer to C-3).** C-3 asks whether **N-a's scope** includes class/scope. P-3
describes what a fold over existing vocabulary would produce. A status-only N-a is compatible with P-1, P-2, P-3 or
P-6; the choice of reading is Q5, not Q4.

---

## 5. What this changes for Q1 and Q5

**C3P-23 (Calculated — Q1's subject is now sharper).** The subject question is not only "which part" but **which
claim**: claim A (derivation) is K-R12's; claim B (transition-per-change) is not. A C-3 answer that denies class/scope
are projected must engage claim A; an answer that denies they have transition kinds engages only claim B.

**C3P-24 (Calculated — Q5's operation identity now has a vocabulary).** The candidate operations are the six readings.
Q5 can be asked as: **which reading, if any, is the intended carrier** of a class/scope change? That is more precise
than "is there an operation," and it does not assume P-1.

**C3P-25 (Calculated — Q5 remains prior to Q4).** "Deferred" and "excluded" are dispositions of a carrier that must
first be identified (C3Q-31). The six readings make the point concrete: if the carrier is P-3 (none), "excluded from
N-a" is coherent but says nothing about change elsewhere; if the carrier is P-4 (another operation), exclusion from
N-a is compatible with change; if the carrier is P-1/P-2/P-6, exclusion from N-a leaves the change in the stream but
outside the fold.

**C3P-26 (Calculated — T-2 and T-4 are refined).** **T-2** now means "separate claim A from claim B, and S-C1 from
S-C2." **T-4** now means "identify which of P-1…P-6 is the intended carrier, or record that none is." The other terms
(T-1, T-3, T-5, T-6) are unchanged (C3Q-29).

**C3P-27 (Calculated — the most important negative result).** **The absence of class/scope transition kinds does not
establish S-C2, and K-R12's projection clause does not establish P-1.** Both directions of the old inference are
blocked. What remains is a genuine question about the carrier and the subject, with a current-constraint finding
(P-3/P-5) that is revisable and not semantic.

---

## 6. What this record does not do

**C3P-28 (Observed).**

| Untouched | State |
|---|---|
| **C-3** | **not answered**; no reading chosen; no request drafted |
| **C-1 / C-2** | determined / non-decision; not reopened |
| **Q1…Q5; T-1…T-6** | open; T-2 and T-4 refined, not answered |
| **N-a / F-1…F-6** | unanswered; none selected, ranked or excluded; **F-6 is not a default** |
| **I-1 / I-2 / I-3** | I-3 excluded by K-R12; the others open |
| **N-b / N-c / N-d** | untouched |
| **O-8a / O-9 / O-10b / O-11** | unchanged; authority-promotion only |
| **DN-6a / DN-6b; DN-1; DQ-4** | unchanged / settled / open |
| **§11's adoption; §8.1's collision** | not decided |
| **map / requirements** | **118**; DQ-3's capture gap acknowledged |

**C3P-29 (Observed — no rule, predicate, schema or capture).** No authority rule, validity predicate, event kind,
field, key, fold, entry architecture or canonical record is proposed or created.

---

## 7. Traceability

**C3P-30 (Observed).**

| What | Where |
|---|---|
| The projection clause | K-R12 (`:111`) |
| Claim A vs claim B | C3P-06/07; C3P-08 |
| The kind list vs "plus … changes" | K-R12 (`:111`); C3P-08 |
| S-C1 / S-C2 | C3Q-19…C3Q-22 |
| The five questions; Q5 prior to Q4 | C3Q-14; C3Q-31 |
| The six readings | C3P-10 |
| Independent rebuild; projection isolation | K-R14 (`:113`); `:367` |
| Cross-stream join is read-time only | K-R13 (`:112`) |
| No cross-stream status conferral | ECI-16 |
| Stored field as cached projection | A2-11; `:366` |
| Formation writes class/scope | `:604`; `:404`; ECI-04/05 |
| No new kinds/fields/requirements | NSD-09; A-6 (NAE-31) |
| Scope vs authority domain | RX-17; C3Q-05 |
| The owner's distinction | 2026-11-20; C3P-01 |

---

## 8. Limits and record status

**C3P-31 (Observed — method).** **Newly read:** K-R12 (`:111`), K-R13 (`:112`), K-R14 (`:113`), `:366`/`:367`,
§6.1 (`:341`), §11 (`:604`), §8.1 (`:404`), A2-11/12, ECI-04/05/11/16, C3Q-14/19…22/29/31, NSD-09, A-6. **Not
inspected:** any implementation's class/scope handling, r32's bytes, the frozen preregistration, v9 beyond §7. No
test, probe or model call was run (C3P-32).

**C3P-32 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The record is textual and definitional.

**C3P-33 (Observed — the labels' limits).** C3P-01…C3P-05, C3P-08, C3P-17, C3P-28…C3P-30 are **Observed** or
**Source-derived**. C3P-06/07, C3P-09…C3P-16, C3P-18…C3P-27 are **Calculated** except where marked **Narrowed** or
**Unestablished**. Nothing is offered as a provision of r33 beyond the quotations already recorded.

**C3P-34 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| C-3 answered; reading chosen; request drafted | **0** |
| Fold, entry model, domain, event kind, field or schema selected | **0** |
| C-1 reopened; C-2 made a decision | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** **The projection clause entails derivation (claim A), not a transition kind per changed value
(claim B)** (C3P-06/07), and its own wording separates the kind list from "plus … changes" (C3P-08). **Six readings**
carry a class/scope change differently: own kind (P-1), generic entry (P-2), static after formation (P-3), multi-source
join (P-4), stored cache (P-5), piggyback (P-6) (C3P-10). Under the **current** constraint set only **P-3/P-5** are
representable, so class/scope would be static post-formation **in any representation buildable today** — a
constraint finding, **not** a semantic prohibition, and **not** an answer to C-3 (C3P-19…C3P-22). **Q1 now asks which
claim is at issue; Q5 asks which reading is the carrier; Q5 remains prior to Q4** (C3P-23…C3P-27). **Nothing is
answered, chosen, requested or authorised.**
